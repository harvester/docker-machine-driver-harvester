package harvester

import (
	"fmt"
	"strings"

	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func (d *Driver) PreCreateCheck() error {
	// server version
	serverVersion, err := d.getSetting("server-version")
	if err != nil {
		return err
	}
	d.ServerVersion = serverVersion.Value
	if strings.HasPrefix(d.ServerVersion, "v0.1.0") {
		return fmt.Errorf("current harvester server version is %s, only support v0.2.0+", d.ServerVersion)
	}

	// vm already exists
	if _, err = d.getVM(); err == nil {
		return fmt.Errorf("machine %s already exists in namespace %s", d.MachineName, d.VMNamespace)
	}

	// keypair check
	if d.KeyPairName != "" {
		keypair, err := d.getKeyPair(d.KeyPairName)
		if err != nil {
			return err
		}

		// keypair validated
		keypairValidated := false
		for _, condition := range keypair.Status.Conditions {
			if condition.Type == harvsterv1.KeyPairValidated && condition.Status == corev1.ConditionTrue {
				keypairValidated = true
				break
			}
		}
		if !keypairValidated {
			return fmt.Errorf("keypair %s is not validated", keypair.Name)
		}

		d.SSHPublicKey = keypair.Spec.PublicKey
	}

	// image and storageClass check
	if d.DiskInfo != nil {
		for _, disk := range d.DiskInfo.Disks {
			if disk.ImageName != "" {
				if _, err = d.getImage(disk.ImageName); err != nil {
					return err
				}
			}
			if disk.StorageClassName != "" {
				if _, err = d.getStorageClass(disk.StorageClassName); err != nil {
					return err
				}
			}
		}
	} else {
		// Compatible with older versions
		if _, err = d.getImage(d.ImageName); err != nil {
			return err
		}
	}

	// network check
	if d.NetworkInfo != nil {
		for _, networkInterface := range d.NetworkInfo.NetworkInterfaces {
			if _, err = d.getNetwork(networkInterface.NetworkName); err != nil {
				return err
			}
		}
	} else {
		// Compatible with older versions
		if d.NetworkName != "" {
			if _, err = d.getNetwork(d.NetworkName); err != nil {
				return err
			}
		}
	}

	return nil
}
