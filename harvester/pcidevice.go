package harvester

import (
	"encoding/json"
	"fmt"
)

type HostDeviceInfo struct {
	HostDeviceRequests []HostDeviceRequest `json:"hostDeviceRequests"`
}

type HostDeviceRequest struct {
	Name       string `json:"name"`
	DeviceName string `json:"deviceName"`
}

func parseHostDeviceInfo(info string) (*HostDeviceInfo, error) {
	hdi := &HostDeviceInfo{}

	if err := json.Unmarshal([]byte(info), hdi); err != nil {
		return nil, fmt.Errorf("error unmarshalling host device information string: %w", err)
	}
	return hdi, nil
}
