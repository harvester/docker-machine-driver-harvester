package harvester

import (
	"encoding/json"
	"fmt"

	kubevirtv1 "kubevirt.io/api/core/v1"
)

type HostDeviceInfo struct {
	HostDevices []kubevirtv1.HostDevice `json:"hostDevices"`
}

func parseHostDeviceInfo(info string) (*HostDeviceInfo, error) {
	hdi := &HostDeviceInfo{}

	if err := json.Unmarshal([]byte(info), hdi); err != nil {
		return nil, fmt.Errorf("error unmarshalling host device information string: %w", err)
	}
	return hdi, nil
}
