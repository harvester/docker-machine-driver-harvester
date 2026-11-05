package harvester

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "go.yaml.in/yaml/v3"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const testWindowsCloudInitFilePath = "testdata/cloudinit_rke2_windows.yaml"

func TestCheckOS(t *testing.T) {
	for _, os := range []string{"", "linux", "windows", "Windows"} {
		assert.NoError(t, checkOS(os), os)
	}
	assert.Error(t, checkOS("freebsd"))
}

func TestMergeCloudInitWindows(t *testing.T) {
	d := NewDriver("m", "")
	d.OS = osWindows
	d.CloudConfig = testWindowsCloudInitFilePath
	d.SSHPublicKey = "ssh-ed25519 AAAA test"
	d.SSHPassword = "Passw0rd!"
	d.UserData = `## template: jinja
#cloud-config
runcmd:
- powershell -Command "Get-Date"
`
	got, _, err := d.mergeCloudInit()
	require.NoError(t, err)

	// Cloudbase-Init needs "#cloud-config" on the first line.
	assert.True(t, strings.HasPrefix(got, "#cloud-config\n"), got)

	cf := map[string]any{}
	require.NoError(t, yaml.Unmarshal([]byte(got), &cf))
	// cloud-init only keys are dropped for Windows
	for _, k := range []string{"ssh_authorized_keys", "chpasswd", "ssh_pwauth", "password", "user"} {
		assert.NotContains(t, cf, k)
	}
	// no user without a user name
	assert.NotContains(t, cf, "users")
	// user runcmd first, then the Rancher install script
	assert.Equal(t, []any{`powershell -Command "Get-Date"`, `powershell C:\install.ps1`}, cf["runcmd"])
	assert.Equal(t, "wincl-pool-x7k2q", cf["set_hostname"])
	assert.Equal(t, "wincl-pool-x7k2q", windowsHostname(got))
}

func TestMergeCloudInitWindowsUser(t *testing.T) {
	d := NewDriver("m", "")
	d.OS = osWindows
	d.SSHUser = "rancher"
	d.SSHPassword = "Passw0rd!"
	got, _, err := d.mergeCloudInit()
	require.NoError(t, err)
	cf := map[string]any{}
	require.NoError(t, yaml.Unmarshal([]byte(got), &cf))
	assert.Equal(t, []any{map[string]any{"name": "rancher", "passwd": "Passw0rd!", "primary_group": "Administrators"}}, cf["users"])
}

func TestMergeCloudInitLinuxUnchanged(t *testing.T) {
	d := NewDriver("m", "")
	d.SSHUser = "ubuntu"
	d.SSHPassword = "pw"
	d.UserData = "## template: jinja\n#cloud-config\n"
	got, _, err := d.mergeCloudInit()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(got, "## template: jinja\n#cloud-config\n"), got)
	assert.Contains(t, got, "chpasswd")
}

func TestWindowsHostname(t *testing.T) {
	assert.Equal(t, "abc-123", windowsHostname("#cloud-config\nset_hostname: ABC-123\n"))
	assert.Equal(t, "", windowsHostname("#cloud-config\nhostname: abc\n"))
	assert.Equal(t, "", windowsHostname("#cloud-config\nset_hostname: not_a.label\n"))
	assert.Equal(t, "", windowsHostname(":::"))
}

func TestConfigureWindows(t *testing.T) {
	vm := &kubevirtv1.VirtualMachine{Spec: kubevirtv1.VirtualMachineSpec{Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{}}}
	smm := true
	vm.Spec.Template.Spec.Domain.Features = &kubevirtv1.Features{SMM: &kubevirtv1.FeatureState{Enabled: &smm}}

	(&Driver{OS: osLinux}).ConfigureWindows(vm, "h")
	assert.Nil(t, vm.Spec.Template.Spec.Domain.Clock, "linux VMs are untouched")

	d := &Driver{OS: osWindows}
	d.ConfigureWindows(vm, "wincl-pool-x7k2q")
	d.ConfigureWindows(vm, "wincl-pool-x7k2q") // idempotent
	spec := vm.Spec.Template.Spec
	assert.Equal(t, osWindows, vm.Labels[osLabelKey])
	assert.Equal(t, osWindows, vm.Spec.Template.ObjectMeta.Labels[osLabelKey])
	assert.Equal(t, "wincl-pool-x7k2q", spec.Hostname)
	require.NotNil(t, spec.Domain.Features.Hyperv)
	assert.True(t, *spec.Domain.Features.Hyperv.Relaxed.Enabled)
	assert.Equal(t, uint32(8191), *spec.Domain.Features.Hyperv.Spinlocks.Retries)
	assert.True(t, *spec.Domain.Features.Hyperv.SyNICTimer.Direct.Enabled)
	assert.True(t, *spec.Domain.Features.SMM.Enabled, "EFI/SMM settings are kept")
	assert.True(t, *spec.Domain.Clock.Timer.Hyperv.Enabled)
	assert.False(t, *spec.Domain.Clock.Timer.HPET.Enabled)
	assert.NotNil(t, spec.Domain.Clock.UTC)
	assert.Len(t, spec.Domain.Devices.Inputs, 1)
}
