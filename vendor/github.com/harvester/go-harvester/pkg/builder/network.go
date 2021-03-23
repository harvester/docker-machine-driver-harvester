package builder

import (
	"fmt"

	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

const (
	defaultVMManagementNetworkName   = "default"
	defaultVMManagementInterfaceName = "default"
	defaultVMInterfaceModel          = "virtio"
)

func (v *VMBuilder) generateNICName() string {
	return fmt.Sprintf("nic-%d", len(v.nicNames))
}

func (v *VMBuilder) ManagementNetwork(bridge bool) *VMBuilder {
	// Networks
	networks := v.vm.Spec.Template.Spec.Networks
	networks = append(networks, kubevirtv1.Network{
		Name: defaultVMManagementNetworkName,
		NetworkSource: kubevirtv1.NetworkSource{
			Pod: &kubevirtv1.PodNetwork{},
		},
	})
	v.vm.Spec.Template.Spec.Networks = networks
	// Interfaces
	interfaces := v.vm.Spec.Template.Spec.Domain.Devices.Interfaces
	nic := kubevirtv1.Interface{
		Name:  defaultVMManagementInterfaceName,
		Model: defaultVMInterfaceModel,
	}
	if bridge {
		nic.InterfaceBindingMethod = kubevirtv1.InterfaceBindingMethod{
			Bridge: &kubevirtv1.InterfaceBridge{},
		}
	} else {
		nic.InterfaceBindingMethod = kubevirtv1.InterfaceBindingMethod{
			Masquerade: &kubevirtv1.InterfaceMasquerade{},
		}
	}
	interfaces = append(interfaces, nic)

	v.vm.Spec.Template.Spec.Domain.Devices.Interfaces = interfaces
	return v
}

func (v *VMBuilder) Bridge(networkName, networkModel string) *VMBuilder {
	nicName := v.generateNICName()
	v.nicNames = append(v.nicNames, nicName)
	// Networks
	networks := v.vm.Spec.Template.Spec.Networks
	networks = append(networks, kubevirtv1.Network{
		Name: nicName,
		NetworkSource: kubevirtv1.NetworkSource{
			Multus: &kubevirtv1.MultusNetwork{
				NetworkName: networkName,
				Default:     false,
			},
		},
	})
	v.vm.Spec.Template.Spec.Networks = networks
	// Interfaces
	interfaces := v.vm.Spec.Template.Spec.Domain.Devices.Interfaces
	interfaces = append(interfaces, kubevirtv1.Interface{
		Name:  nicName,
		Model: networkModel,
		InterfaceBindingMethod: kubevirtv1.InterfaceBindingMethod{
			Bridge: &kubevirtv1.InterfaceBridge{},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Interfaces = interfaces
	return v
}
