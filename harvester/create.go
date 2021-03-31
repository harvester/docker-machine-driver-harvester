package harvester

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/harvester/go-harvester/pkg/builder"
	goharverrors "github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/mcnutils"
	"github.com/rancher/machine/libmachine/ssh"
	"github.com/rancher/machine/libmachine/state"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

const (
	keypairNamespace = "harvester-system"
)

func (d *Driver) PreCreateCheck() error {
	c, err := d.getClient()
	if err != nil {
		return err
	}

	// vm doesn't exist
	_, err = c.VirtualMachines.Get(d.Namespace, d.MachineName)
	if err == nil {
		return fmt.Errorf("machine %s already exists", d.MachineName)
	}

	// image exist
	image, err := c.Images.Get(d.Namespace, d.ImageName)
	if err != nil {
		if goharverrors.IsNotFound(err) {
			return fmt.Errorf("image %s doesn't exist in namespace %s", d.ImageName, d.Namespace)
		}
		return err
	}

	// image succeed
	if image.Status.Progress != 100 {
		return fmt.Errorf("image %s's progress %d != 100", image.Name, image.Status.Progress)
	}
	d.ImageDownloadURL = image.Status.DownloadURL

	if d.KeyPairName != "" {
		keypair, err := c.Keypairs.Get(keypairNamespace, d.KeyPairName)
		if err != nil {
			if goharverrors.IsNotFound(err) {
				return fmt.Errorf("keypair %s doesn't exist in namespace %s", d.KeyPairName, keypairNamespace)
			}
			return err
		}

		// keypair validated
		keypairValidated := false
		for _, condition := range keypair.Status.Conditions {
			if condition.Type == "validated" && condition.Status == "True" {
				keypairValidated = true
			}
		}
		if !keypairValidated {
			return fmt.Errorf("keypair %s is not validated", keypair.Name)
		}

		d.SSHPublicKey = keypair.Spec.PublicKey
	}

	// network exist
	if d.NetworkType != networkTypePod {
		_, err = c.Networks.Get(d.Namespace, d.NetworkName)
		if err != nil {
			if goharverrors.IsNotFound(err) {
				return fmt.Errorf("network %s doesn't exist in namespace %s", d.KeyPairName, d.Namespace)
			}
			return err
		}
	}

	return err
}

func (d *Driver) Create() error {
	c, err := d.getClient()
	if err != nil {
		return err
	}

	if err = d.createKeyPair(); err != nil {
		return err
	}

	userData, networkData := d.createCloudInit()

	var dataVolumeOption *builder.DataVolumeOption
	serverVersion, err := c.Settings.Get("server-version")
	if err != nil {
		return err
	}
	supportLiveMigrate := !strings.HasPrefix(serverVersion.Value, "v0.1.0")
	if supportLiveMigrate {
		dataVolumeOption = &builder.DataVolumeOption{
			VolumeMode:       corev1.PersistentVolumeBlock,
			AccessMode:       corev1.ReadWriteMany,
			StorageClassName: pointer.StringPtr("longhorn-" + d.ImageName),
		}
	} else {
		dataVolumeOption = &builder.DataVolumeOption{
			HTTPURL:    d.ImageDownloadURL,
			VolumeMode: corev1.PersistentVolumeFilesystem,
			AccessMode: corev1.ReadWriteOnce,
		}
	}
	dataVolumeOption.ImageID = fmt.Sprintf("%s/%s", d.Namespace, d.ImageName)
	// create vm
	vmBuilder := builder.NewVMBuilder("docker-machine-driver-harvester").
		Namespace(d.Namespace).Name(d.MachineName).
		CPU(d.CPU).Memory(d.MemorySize).
		Image(d.DiskSize, d.DiskBus, dataVolumeOption).
		EvictionStrategy(supportLiveMigrate).
		CloudInit(userData, networkData)

	if d.KeyPairName != "" {
		vmBuilder = vmBuilder.SSHKey(d.KeyPairName)
	}

	if d.NetworkType != networkTypePod {
		vmBuilder = vmBuilder.Bridge(d.NetworkName, d.NetworkModel)
	} else {
		vmBuilder = vmBuilder.ManagementNetwork(true)
	}

	if _, err = c.VirtualMachines.Create(vmBuilder.Run()); err != nil {
		return err
	}

	if err = d.waitForState(state.Running); err != nil {
		return err
	}
	if err = d.waitForIP(); err != nil {
		return err
	}
	ip, err := d.GetIP()
	if err != nil {
		return err
	}
	d.IPAddress = ip
	log.Debugf("Get machine ip: %s", d.IPAddress)
	return nil
}

func (d *Driver) waitForState(desiredState state.State) error {
	log.Debugf("Waiting for node become %s", desiredState)
	if err := mcnutils.WaitForSpecific(drivers.MachineInState(d, desiredState), 120, 5*time.Second); err != nil {
		return fmt.Errorf("Too many retries waiting for machine to be %s.  Last error: %s", desiredState, err)
	}
	return nil
}

func (d *Driver) waitForIP() error {
	ipIsNotEmpty := func() bool {
		ip, _ := d.GetIP()
		return ip != ""
	}
	log.Debugf("Waiting for node get ip")
	if err := mcnutils.WaitForSpecific(ipIsNotEmpty, 120, 5*time.Second); err != nil {
		return fmt.Errorf("Too many retries waiting for get machine's ip.  Last error: %s", err)
	}
	return nil
}

func (d *Driver) waitForReady() error {
	if err := d.waitForState(state.Running); err != nil {
		return err
	}
	return d.waitForIP()
}

func (d *Driver) waitForRestart(oldUID string) error {
	restarted := func() bool {
		vmi, err := d.getVMI()
		if err != nil {
			return false
		}
		return oldUID != string(vmi.UID)
	}
	log.Debugf("Waiting for node restarted")
	if err := mcnutils.WaitForSpecific(restarted, 120, 5*time.Second); err != nil {
		return fmt.Errorf("Too many retries waiting for machine restart.  Last error: %s", err)
	}
	return d.waitForReady()
}

func (d *Driver) createKeyPair() error {
	keyPath := d.GetSSHKeyPath()
	publicKeyFile := keyPath + ".pub"
	if d.SSHPrivateKeyPath == "" {
		log.Debugf("Creating New SSH Key")
		if err := ssh.GenerateSSHKey(keyPath); err != nil {
			return err
		}
	} else {
		log.Debugf("Using SSHPrivateKeyPath: %s", d.SSHPrivateKeyPath)
		if err := mcnutils.CopyFile(d.SSHPrivateKeyPath, keyPath); err != nil {
			return err
		}
		if d.KeyPairName != "" {
			log.Debugf("Using existing harvester key pair: %s", d.KeyPairName)
			return nil
		}
		if err := mcnutils.CopyFile(d.SSHPrivateKeyPath+".pub", publicKeyFile); err != nil {
			return err
		}
	}

	publicKey, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}
	log.Debugf("Using public Key: %s", publicKeyFile)
	d.SSHPublicKey = string(publicKey)
	return nil
}
