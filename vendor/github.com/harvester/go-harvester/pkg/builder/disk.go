package builder

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	cdiv1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func (v *VMBuilder) generateDiskName() string {
	return fmt.Sprintf("disk-%d", len(v.dataVolumeNames))
}

func (v *VMBuilder) Blank(diskSize, diskBus string) *VMBuilder {
	return v.DataVolume(diskSize, diskBus, nil)
}

func (v *VMBuilder) Image(diskSize, diskBus string, opt *DataVolumeOption) *VMBuilder {
	return v.DataVolume(diskSize, diskBus, opt)
}

func (v *VMBuilder) SSHKey(sshKeyName string) *VMBuilder {
	v.sshNames = append(v.sshNames, sshKeyName)
	return v
}

type DataVolumeOption struct {
	HTTPURL          string
	VolumeMode       corev1.PersistentVolumeMode
	AccessMode       corev1.PersistentVolumeAccessMode
	StorageClassName *string
}

func (v *VMBuilder) DataVolume(diskSize, diskBus string, opt *DataVolumeOption) *VMBuilder {
	if opt == nil {
		opt = &DataVolumeOption{
			VolumeMode: corev1.PersistentVolumeBlock,
			AccessMode: corev1.ReadWriteMany,
		}
	}
	diskName := v.generateDiskName()
	dataVolumeName := fmt.Sprintf("%s-%s-%s", v.vm.Name, diskName, rand.String(5))
	v.dataVolumeNames = append(v.dataVolumeNames, dataVolumeName)
	// DataVolumeTemplates
	dataVolumeTemplates := v.vm.Spec.DataVolumeTemplates
	dataVolumeSpecSource := cdiv1alpha1.DataVolumeSource{
		Blank: &cdiv1alpha1.DataVolumeBlankImage{},
	}

	if opt.HTTPURL != "" {
		dataVolumeSpecSource = cdiv1alpha1.DataVolumeSource{
			HTTP: &cdiv1alpha1.DataVolumeSourceHTTP{
				URL: opt.HTTPURL,
			},
		}
	}
	dataVolumeTemplate := kubevirtv1.DataVolumeTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:        dataVolumeName,
			Labels:      nil,
			Annotations: nil,
		},
		Spec: cdiv1alpha1.DataVolumeSpec{
			Source: dataVolumeSpecSource,
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					opt.AccessMode,
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(diskSize),
					},
				},
				VolumeMode:       &opt.VolumeMode,
				StorageClassName: opt.StorageClassName,
			},
		},
	}
	dataVolumeTemplates = append(dataVolumeTemplates, dataVolumeTemplate)
	v.vm.Spec.DataVolumeTemplates = dataVolumeTemplates
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	disks = append(disks, kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			DataVolume: &kubevirtv1.DataVolumeSource{
				Name: dataVolumeName,
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}

func (v *VMBuilder) ExistingDataVolume(dataVolumeName, diskBus string) *VMBuilder {
	diskName := v.generateDiskName()
	v.dataVolumeNames = append(v.dataVolumeNames, dataVolumeName)
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	disks = append(disks, kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			DataVolume: &kubevirtv1.DataVolumeSource{
				Name: dataVolumeName,
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}

func (v *VMBuilder) ContainerDisk(diskName, diskBus, imageName, ImagePullPolicy string, isCDRom bool) *VMBuilder {
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	diskDevice := kubevirtv1.DiskDevice{
		Disk: &kubevirtv1.DiskTarget{
			Bus: diskBus,
		},
	}
	if isCDRom {
		diskDevice = kubevirtv1.DiskDevice{
			CDRom: &kubevirtv1.CDRomTarget{
				Bus: diskBus,
			},
		}
	}
	disks = append(disks, kubevirtv1.Disk{
		Name:       diskName,
		DiskDevice: diskDevice,
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			ContainerDisk: &kubevirtv1.ContainerDiskSource{
				Image:           imageName,
				ImagePullPolicy: corev1.PullPolicy(ImagePullPolicy),
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}

func (v *VMBuilder) Container(diskName, diskBus, imageName, ImagePullPolicy string) *VMBuilder {
	return v.ContainerDisk(diskName, diskBus, imageName, ImagePullPolicy, false)
}

func (v *VMBuilder) CDRom(diskName, diskBus, imageName, ImagePullPolicy string) *VMBuilder {
	return v.ContainerDisk(diskName, diskBus, imageName, ImagePullPolicy, true)
}
