package builder

import (
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

func (v *VMBuilder) CloudInit(userData, networkData string) *VMBuilder {
	diskName := "cloudinitdisk"
	diskBus := "virtio"
	// Disks
	var (
		diskExist bool
		diskIndex int
	)
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	for i, disk := range disks {
		if disk.Name == diskName {
			diskExist = true
			diskIndex = i
			break
		}
	}

	disk := kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	}
	if diskExist {
		disks[diskIndex] = disk
	} else {
		disks = append(disks, disk)
	}

	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks

	// Volumes
	var (
		volumeExist bool
		volumeIndex int
	)
	volumes := v.vm.Spec.Template.Spec.Volumes
	for i, volume := range volumes {
		if volume.Name == diskName {
			volumeExist = true
			volumeIndex = i
			break
		}
	}
	volume := kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			CloudInitNoCloud: &kubevirtv1.CloudInitNoCloudSource{
				UserData:    userData,
				NetworkData: networkData,
			},
		},
	}
	if volumeExist {
		volumes[volumeIndex] = volume
	} else {
		volumes = append(volumes, volume)
	}
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}
