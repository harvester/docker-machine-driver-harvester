package harvester

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func TestGuestHostname(t *testing.T) {
	assert.Equal(t, "win-123", guestHostname("#cloud-config\nset_hostname: WIN-123\n"))
	assert.Equal(t, "lin-123", guestHostname("#cloud-config\nhostname: lin-123\n"))
	assert.Equal(t, "", guestHostname("#cloud-config\nhostname: not_a.label\n"))
	assert.Equal(t, "", guestHostname("#cloud-config\nruncmd: []\n"))
	assert.Equal(t, "", guestHostname(":::"))
}

func TestGuestNodeName(t *testing.T) {
	d := NewDriver("Pool-ABC-x7k2q", "")
	assert.Equal(t, "short-1", d.guestNodeName("", "#cloud-config\nhostname: short-1\n"), "secret user data is used when inline has no hostname")
	assert.Equal(t, "inline-1", d.guestNodeName("#cloud-config\nset_hostname: inline-1\n", "#cloud-config\nhostname: short-1\n"))
	assert.Equal(t, "pool-abc-x7k2q", d.guestNodeName("", ""), "falls back to the machine name")

	d = NewDriver("bad_name", "")
	assert.Equal(t, "", d.guestNodeName(""))
}

func TestConfigureGuestNodeName(t *testing.T) {
	vm := &kubevirtv1.VirtualMachine{Spec: kubevirtv1.VirtualMachineSpec{Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{}}}
	ConfigureGuestNodeName(vm, "")
	assert.Empty(t, vm.Labels)

	ConfigureGuestNodeName(vm, "wintest2-90cdf2")
	assert.Equal(t, "wintest2-90cdf2", vm.Labels[guestNodeNameLabelKey])
	assert.Equal(t, "wintest2-90cdf2", vm.Spec.Template.ObjectMeta.Labels[guestNodeNameLabelKey])
}
