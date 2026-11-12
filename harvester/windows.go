package harvester

import (
	"fmt"
	"strings"

	yaml "go.yaml.in/yaml/v3"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const (
	osLinux   = "linux"
	osWindows = "windows"

	// osLabelKey is the label the Harvester UI uses to tell Windows VMs apart.
	osLabelKey = "harvesterhci.io/os"

	// userDataWindowsUserTemplate creates a local administrator through
	// Cloudbase-Init's cloud-config "users" plugin. Cloudbase-Init does not
	// understand cloud-init's user/password/chpasswd/ssh_authorized_keys keys.
	userDataWindowsUserTemplate = `
users:
- name: %s
  passwd: %s
  primary_group: Administrators`

	// cloudbaseInitHostnameKey is the cloud-config key Cloudbase-Init reads
	// the hostname from; rancher-machine writes it for Windows machines.
	cloudbaseInitHostnameKey = "set_hostname"
)

func checkOS(os string) error {
	switch strings.ToLower(os) {
	case "", osLinux, osWindows:
		return nil
	default:
		return fmt.Errorf("unsupported harvester os %q, must be %s or %s", os, osLinux, osWindows)
	}
}

func (d *Driver) isWindows() bool {
	return strings.EqualFold(d.OS, osWindows)
}

// windowsHostname returns the hostname Rancher asked for in the Cloudbase-Init
// cloud-config, if it is a valid DNS label.
//
// The driver copies it to the VMI's spec.hostname so that the NoCloud
// meta-data local-hostname matches. Otherwise Cloudbase-Init's SetHostName
// plugin sets the computer name from the VM name (truncated to 15 NetBIOS
// characters), which collides across machines of the same pool and differs
// from the name Rancher expects.
func windowsHostname(userData string) string {
	cf := map[string]any{}
	if err := yaml.Unmarshal([]byte(userData), &cf); err != nil {
		return ""
	}
	hostname, ok := cf[cloudbaseInitHostnameKey].(string)
	if !ok {
		return ""
	}
	return validNodeName(hostname)
}

// ConfigureWindows applies the settings Harvester's built-in Windows VM
// templates use: the OS label, Hyper-V enlightenments, a Hyper-V clock and a
// USB tablet for the console pointer.
func (d *Driver) ConfigureWindows(vm *kubevirtv1.VirtualMachine, hostname string) {
	if !d.isWindows() {
		return
	}
	enabled := func() *kubevirtv1.FeatureState {
		t := true
		return &kubevirtv1.FeatureState{Enabled: &t}
	}
	spinlocks := uint32(8191)
	present := true
	notPresent := false

	if vm.Labels == nil {
		vm.Labels = map[string]string{}
	}
	vm.Labels[osLabelKey] = osWindows
	if vm.Spec.Template.ObjectMeta.Labels == nil {
		vm.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}
	vm.Spec.Template.ObjectMeta.Labels[osLabelKey] = osWindows

	spec := &vm.Spec.Template.Spec
	if hostname != "" {
		spec.Hostname = hostname
	}

	if spec.Domain.Features == nil {
		spec.Domain.Features = &kubevirtv1.Features{}
	}
	features := spec.Domain.Features
	features.ACPI = kubevirtv1.FeatureState{Enabled: &present}
	features.APIC = &kubevirtv1.FeatureAPIC{FeatureState: *enabled()}
	features.Hyperv = &kubevirtv1.FeatureHyperv{
		Relaxed:    enabled(),
		VAPIC:      enabled(),
		Spinlocks:  &kubevirtv1.FeatureSpinlocks{FeatureState: *enabled(), Retries: &spinlocks},
		VPIndex:    enabled(),
		SyNIC:      enabled(),
		SyNICTimer: &kubevirtv1.SyNICTimer{FeatureState: *enabled(), Direct: enabled()},
		IPI:        enabled(),
		Runtime:    enabled(),
		Reset:      enabled(),
	}

	spec.Domain.Clock = &kubevirtv1.Clock{
		ClockOffset: kubevirtv1.ClockOffset{UTC: &kubevirtv1.ClockOffsetUTC{}},
		Timer: &kubevirtv1.Timer{
			HPET:   &kubevirtv1.HPETTimer{Enabled: &notPresent},
			Hyperv: &kubevirtv1.HypervTimer{Enabled: &present},
			PIT:    &kubevirtv1.PITTimer{TickPolicy: kubevirtv1.PITTickPolicyDelay},
			RTC:    &kubevirtv1.RTCTimer{TickPolicy: kubevirtv1.RTCTickPolicyCatchup},
		},
	}

	hasTablet := false
	for _, input := range spec.Domain.Devices.Inputs {
		if input.Type == kubevirtv1.InputTypeTablet {
			hasTablet = true
		}
	}
	if !hasTablet {
		spec.Domain.Devices.Inputs = append(spec.Domain.Devices.Inputs, kubevirtv1.Input{
			Name: "tablet",
			Type: kubevirtv1.InputTypeTablet,
			Bus:  kubevirtv1.InputBusUSB,
		})
	}
}
