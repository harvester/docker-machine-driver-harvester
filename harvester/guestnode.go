package harvester

import (
	"strings"

	yaml "go.yaml.in/yaml/v3"
	"k8s.io/apimachinery/pkg/util/validation"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const (
	// guestNodeNameLabelKey records the Kubernetes node name the guest
	// registers with. The Harvester cloud provider maps nodes to VMs with
	// this label, so it does not depend on the VM name or on the hostname
	// the guest agent reports.
	guestNodeNameLabelKey = "guestcluster.harvesterhci.io/node-name"

	// cloudInitHostnameKey is the cloud-config key rancher-machine writes
	// the hostname to for Linux machines (set_hostname for Windows).
	cloudInitHostnameKey = "hostname"
)

// guestHostname returns the hostname that rancher-machine asked for in the
// cloud-config: set_hostname (Cloudbase-Init) or hostname (cloud-init),
// lowercased, if it is a valid DNS label.
func guestHostname(userData string) string {
	cf := map[string]any{}
	if err := yaml.Unmarshal([]byte(userData), &cf); err != nil {
		return ""
	}
	for _, key := range []string{cloudbaseInitHostnameKey, cloudInitHostnameKey} {
		if hostname, ok := cf[key].(string); ok {
			return validNodeName(hostname)
		}
	}
	return ""
}

// validNodeName lowercases name and returns it if it is a valid DNS label,
// which also makes it a valid label value.
func validNodeName(name string) string {
	name = strings.ToLower(name)
	if len(validation.IsDNS1123Label(name)) > 0 {
		return ""
	}
	return name
}

// guestNodeName returns the node name the guest will register with: the
// hostname from the user data, or the machine name when there is none.
func (d *Driver) guestNodeName(userDataSources ...string) string {
	for _, userData := range userDataSources {
		if hostname := guestHostname(userData); hostname != "" {
			return hostname
		}
	}
	return validNodeName(d.MachineName)
}

// ConfigureGuestNodeName labels the VM and its VMI template with the guest
// node name. An empty name leaves the VM unchanged.
func ConfigureGuestNodeName(vm *kubevirtv1.VirtualMachine, nodeName string) {
	if nodeName == "" {
		return
	}
	if vm.Labels == nil {
		vm.Labels = map[string]string{}
	}
	vm.Labels[guestNodeNameLabelKey] = nodeName
	if vm.Spec.Template == nil {
		return
	}
	if vm.Spec.Template.ObjectMeta.Labels == nil {
		vm.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}
	vm.Spec.Template.ObjectMeta.Labels[guestNodeNameLabelKey] = nodeName
}
