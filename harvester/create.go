package harvester

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/harvester/harvester/pkg/builder"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/mcnutils"
	"github.com/rancher/machine/libmachine/ssh"
	"github.com/rancher/machine/libmachine/state"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

const (
	rootDiskName  = "disk-0"
	interfaceName = "nic-0"
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

	// vm doesn't exist
	if _, err = d.getVM(); err == nil {
		return fmt.Errorf("machine %s already exists in namespace %s", d.MachineName, d.VMNamespace)
	}

	// image exist
	if _, err = d.getImage(); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("image %s doesn't exist", d.ImageName)
		}
		return err
	}

	if d.KeyPairName != "" {
		keypair, err := d.getKeyPair()
		if err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("keypair %s doesn't exist", d.KeyPairName)
			}
			return err
		}

		// keypair validated
		keypairValidated := false
		for _, condition := range keypair.Status.Conditions {
			if condition.Type == harvsterv1.KeyPairValidated && condition.Status == corev1.ConditionTrue {
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
		if _, err = d.getNetwork(); err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("network %s doesn't exist", d.NetworkName)
			}
			return err
		}
	}

	return err
}

func (d *Driver) Create() error {
	// create keypair
	if err := d.createKeyPair(); err != nil {
		return err
	}
	// create vm
	cloudInitSource, cloudConfigSecret, err := d.buildCloudInit()
	if err != nil {
		return err
	}
	imageNamespace, imageName, err := NamespacedNamePartsByDefault(d.ImageName, d.VMNamespace)
	if err != nil {
		return err
	}
	pvcOption := &builder.PersistentVolumeClaimOption{
		ImageID:          fmt.Sprintf("%s/%s", imageNamespace, imageName),
		VolumeMode:       corev1.PersistentVolumeBlock,
		AccessMode:       corev1.ReadWriteMany,
		StorageClassName: pointer.StringPtr(builder.BuildImageStorageClassName("", imageName)),
	}
	vmBuilder := builder.NewVMBuilder("docker-machine-driver-harvester").
		Namespace(d.VMNamespace).Name(d.MachineName).CPU(d.CPU).Memory(d.MemorySize).
		PVCDisk(rootDiskName, builder.DiskBusVirtio, false, false, 1, d.DiskSize, "", pvcOption).
		CloudInitDisk(builder.CloudInitDiskName, builder.DiskBusVirtio, false, 0, *cloudInitSource).
		EvictionStrategy(true).DefaultPodAntiAffinity().Run(false)

	if d.KeyPairName != "" {
		vmBuilder = vmBuilder.SSHKey(d.KeyPairName)
	}
	interfaceType := builder.NetworkInterfaceTypeBridge
	networkName := d.NetworkName
	if d.NetworkType == networkTypePod {
		networkName = ""
	}
	vm, err := vmBuilder.NetworkInterface(interfaceName, d.NetworkModel, "", interfaceType, networkName).VM()
	if err != nil {
		return err
	}
	vm.Kind = kubevirtv1.VirtualMachineGroupVersionKind.Kind
	vm.APIVersion = kubevirtv1.GroupVersion.String()
	createdVM, err := d.createVM(vm)
	if err != nil {
		return err
	}
	// create secret
	if cloudConfigSecret != nil {
		cloudConfigSecret.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: vm.APIVersion,
				Kind:       vm.Kind,
				Name:       vm.Name,
				UID:        createdVM.UID,
			},
		}
		if _, err = d.createSecret(cloudConfigSecret); err != nil {
			return err
		}
	}
	// start vm
	if err = d.Start(); err != nil {
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
