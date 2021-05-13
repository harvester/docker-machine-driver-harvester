package harvester

import (
	"errors"
	"fmt"

	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
)

const (
	defaultNamespace = "default"

	defaultInClusterHost = "harvester.harvester-system"
	defaultInClusterPort = 8443

	defaultCPU          = 2
	defaultMemorySize   = 4
	defaultDiskSize     = 40
	defaultDiskBus      = "virtio"
	defaultNetworkModel = "virtio"
	networkTypePod      = "pod"
	networkTypeDHCP     = "dhcp"
)

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_HOST",
			Name:   "harvester-host",
			Usage:  "harvester host",
			Value:  defaultInClusterHost,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_PORT",
			Name:   "harvester-port",
			Usage:  "harvester port",
			Value:  defaultInClusterPort,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_USERNAME",
			Name:   "harvester-username",
			Usage:  "harvester username",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_PASSWORD",
			Name:   "harvester-password",
			Usage:  "harvester password",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_CLUSTER_TYPE",
			Name:   "harvester-cluster-type",
			Usage:  "harvester cluster type",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NAMESPACE",
			Name:   "harvester-namespace",
			Usage:  "harvester namespace",
			Value:  defaultNamespace,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_CPU_COUNT",
			Name:   "harvester-cpu-count",
			Usage:  "number of CPUs for the machine",
			Value:  defaultCPU,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_MEMORY_SIZE",
			Name:   "harvester-memory-size",
			Usage:  "size of memory for machine (in GiB)",
			Value:  defaultMemorySize,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_DISK_SIZE",
			Name:   "harvester-disk-size",
			Usage:  "size of disk for machine (in GiB)",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_DISK_BUS",
			Name:   "harvester-disk-bus",
			Usage:  "bus of disk for machine",
			Value:  defaultDiskBus,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_IMAGE_NAME",
			Name:   "harvester-image-name",
			Usage:  "harvester image name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_USER",
			Name:   "harvester-ssh-user",
			Usage:  "SSH username",
			Value:  drivers.DefaultSSHUser,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_SSH_PORT",
			Name:   "harvester-ssh-port",
			Usage:  "SSH port",
			Value:  drivers.DefaultSSHPort,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_PASSWORD",
			Name:   "harvester-ssh-password",
			Usage:  "SSH password",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_KEY_PAIR_NAME",
			Name:   "harvester-key-pair-name",
			Usage:  "harvester key pair name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_PRIVATE_KEY_PATH",
			Name:   "harvester-ssh-private-key-path",
			Usage:  "SSH private key path ",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_TYPE",
			Name:   "harvester-network-type",
			Usage:  "harvester network type",
			Value:  networkTypeDHCP,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_NAME",
			Name:   "harvester-network-name",
			Usage:  "harvester network name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_MODEL",
			Name:   "harvester-network-model",
			Usage:  "harvester network model",
			Value:  defaultNetworkModel,
		},
	}
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Host = flags.String("harvester-host")
	d.Port = flags.Int("harvester-port")
	d.Username = flags.String("harvester-username")
	d.Password = flags.String("harvester-password")
	d.Namespace = flags.String("harvester-namespace")
	d.ClusterType = flags.String("harvester-cluster-type")

	d.CPU = flags.Int("harvester-cpu-count")
	d.MemorySize = fmt.Sprintf("%dGi", flags.Int("harvester-memory-size"))
	d.DiskSize = fmt.Sprintf("%dGi", flags.Int("harvester-disk-size"))
	d.DiskBus = flags.String("harvester-disk-bus")

	d.ImageName = flags.String("harvester-image-name")

	d.SSHUser = flags.String("harvester-ssh-user")
	d.SSHPort = flags.Int("harvester-ssh-port")

	d.KeyPairName = flags.String("harvester-key-pair-name")
	d.SSHPrivateKeyPath = flags.String("harvester-ssh-private-key-path")
	d.SSHPassword = flags.String("harvester-ssh-password")

	d.NetworkType = flags.String("harvester-network-type")

	d.NetworkName = flags.String("harvester-network-name")
	d.NetworkModel = flags.String("harvester-network-model")

	d.SetSwarmConfigFromFlags(flags)

	return d.checkConfig()
}

func (d *Driver) checkConfig() error {
	if d.Username == "" {
		return errors.New("must specify harvester username")
	}
	if d.Password == "" {
		return errors.New("must specify harvester password")
	}
	if d.ImageName == "" {
		return errors.New("must specify harvester image name")
	}
	if d.KeyPairName != "" && d.SSHPrivateKeyPath == "" {
		return errors.New("must specify the ssh private key path of the harvester key pair")
	}
	switch d.NetworkType {
	case networkTypePod:
	case networkTypeDHCP:
		if d.NetworkName == "" {
			return errors.New("must specify harvester network name")
		}
	default:
		return fmt.Errorf("unknown network type %s", d.NetworkType)
	}
	return nil
}
