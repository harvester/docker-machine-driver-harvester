package harvester

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

func UnmarshalDiskInfo(data []byte) (DiskInfo, error) {
	var r DiskInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DiskInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DiskInfo struct {
	Disks []Disk `json:"disks"`
}

type Disk struct {
	ImageName        string `json:"imageName"`
	StorageClassName string `json:"storageClassName"`

	Size      int  `json:"size"`
	BootOrder uint `json:"bootOrder"`

	Bus  string `json:"bus"`
	Type string `json:"type"`

	HotPlugAble bool `json:"hotPlugAble"`
}

func UnmarshalNetworkInfo(data []byte) (NetworkInfo, error) {
	var r NetworkInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NetworkInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NetworkInfo struct {
	NetworkInterfaces []NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct {
	NetworkName string `json:"networkName"`

	MACAddress string `json:"macAddress"`

	Model string `json:"model"`
	Type  string `json:"type"`
}

type VGPUInfo struct {
	VGPURequests []VGPURequest `json:"vGPURequests"`
}

type VGPURequest struct {
	Name       string `json:"name"`
	DeviceName string `json:"deviceName"`
}

func (d *Driver) checkConfig() error {
	if d.KeyPairName != "" && d.SSHPrivateKeyPath == "" {
		return errors.New("must specify the ssh private key path of the harvester key pair")
	}
	if d.DiskInfo != nil {
		for _, disk := range d.DiskInfo.Disks {
			if disk.ImageName == "" && disk.StorageClassName == "" {
				return errors.New("must specify image name or storageClass name in harvester disk info")
			}
			if disk.Size <= 0 {
				return errors.New("must specify disk size in harvester disk info")
			}
		}
	} else {
		// Compatible with older versions
		if d.ImageName == "" {
			return errors.New("must specify harvester image name")
		}
		if d.DiskSize == "0" {
			return errors.New("must specify harvester disk size")
		}
	}
	if d.NetworkInfo != nil {
		for _, networkInterface := range d.NetworkInfo.NetworkInterfaces {
			if networkInterface.NetworkName == "" {
				return errors.New("must specify network name in harvester network info")
			}
		}
	} else {
		// Compatible with older versions
		if d.NetworkName == "" {
			return errors.New("must specify harvester network name")
		}
	}
	return checkNetworkData(d.NetworkData)
}

func checkNetworkData(networkDataStr string) error {
	if networkDataStr == "" {
		return nil
	}

	network, version, err := parserNetworkData(networkDataStr)
	if err != nil {
		return err
	}

	switch version {
	case 1:
		if err = checkNetworkDataV1(network); err != nil {
			return err
		}
	case 2:
		// TODO check network data v2 version format here
	}

	return nil
}

func parserNetworkData(networkDataStr string) (map[string]interface{}, float64, error) {
	var networkData = make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(networkDataStr), &networkData); err != nil {
		return nil, 0, err
	}
	// root section
	var rootSection map[string]interface{}
	networkSection, ok := networkData["network"]
	if ok {
		rootSection = networkSection.(map[string]interface{})
	} else {
		rootSection = networkData
	}

	// network.version
	versionSection, err := mustGetSection(rootSection, "version")
	if err != nil {
		return rootSection, 0, err
	}
	version := versionSection.(float64)
	return rootSection, version, nil
}

func checkNetworkDataV1(network map[string]interface{}) error {
	var defaultGatewayCount, nameServerCount, dhcpAllCount int

	// network.config
	networkConfigSection, err := mustGetSection(network, "config")
	if err != nil {
		return err
	}
	networkConfigs := networkConfigSection.([]interface{})

	for _, networkConfig := range networkConfigs {
		config := networkConfig.(map[string]interface{})
		// network.config[].type
		typeSection, err := mustGetSection(config, "type")
		if err != nil {
			return err
		}
		configType := typeSection.(string)
		switch configType {
		case "physical":
			gatewayCount, dhcpCount, err := getGatewayAndDHCPCount(config)
			if err != nil {
				return err
			}
			defaultGatewayCount += gatewayCount
			nameServerCount += dhcpCount
			dhcpAllCount += dhcpCount
		case "nameserver":
			nameServerAddressCount, err := getNameServerAddressCount(config)
			if err != nil {
				return err
			}
			nameServerCount += nameServerAddressCount
		}
	}

	if defaultGatewayCount > 1 {
		return fmt.Errorf("the number of default gateway cannot greater than 1, but get: %d", defaultGatewayCount)
	}

	if defaultGatewayCount == 0 && dhcpAllCount == 0 {
		return errors.New("static gateway or dhcp is not configured")
	}

	if nameServerCount == 0 {
		return errors.New("nameserver is not configured")
	}

	return nil
}

func getNameServerAddressCount(config map[string]interface{}) (int, error) {
	// network.config[].address
	nameServerAddressesSection, err := mustGetSection(config, "address")
	if err != nil {
		return 0, err
	}
	nameServerAddresses := nameServerAddressesSection.([]interface{})
	return len(nameServerAddresses), nil
}

func getGatewayAndDHCPCount(config map[string]interface{}) (int, int, error) {
	var gatewayCount, dhcpCount int

	// network.config[].subnets
	subnetsSection, err := mustGetSection(config, "subnets")
	if err != nil {
		return 0, 0, err
	}
	networkSubnets := subnetsSection.([]interface{})

	for _, networkSubnet := range networkSubnets {
		subnet := networkSubnet.(map[string]interface{})
		// network.config[].subnets[].type
		subnetTypeSection, err := mustGetSection(subnet, "type")
		if err != nil {
			return 0, 0, err
		}
		subnetType := subnetTypeSection.(string)
		switch subnetType {
		case "dhcp":
			// dhcp not always generate a default route
			// gatewayCount += 1
			dhcpCount += 1
		case "static":
			// network.config[].subnets[].gateway
			gatewaySection := subnet["gateway"]
			if gatewaySection != nil {
				gateway := gatewaySection.(string)
				if gateway != "" {
					gatewayCount += 1
				}
			}
		}
	}
	return gatewayCount, dhcpCount, nil
}

func mustGetSection(m map[string]interface{}, k string) (interface{}, error) {
	section := m[k]
	if section == nil {
		return nil, fmt.Errorf("missing section: %s", k)
	}
	return section, nil
}

func parseVGPUInfo(vGPUInfo string) (*VGPUInfo, error) {
	v := &VGPUInfo{}
	err := json.Unmarshal([]byte(vGPUInfo), v)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling vgpuInfo string")
	}
	return v, nil
}
