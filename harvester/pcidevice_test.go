package harvester

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testcase struct {
	description string
	input       string
	expectation *HostDeviceInfo
}

func Test_parseHostDeviceInfo(t *testing.T) {
	testcases := []testcase{
		{
			description: "empty JSON input",
			input:       `{}`,
			expectation: &HostDeviceInfo{},
		},
		{
			description: "non-empty no device info",
			input:       `{"hostDeviceRequests":[]}`,
			expectation: &HostDeviceInfo{
				HostDeviceRequests: []HostDeviceRequest{},
			},
		},
		{
			description: "single device info",
			input:       `{"hostDeviceRequests":[{"name":"qat","deviceName":"intel.com/qat"}]}`,
			expectation: &HostDeviceInfo{
				HostDeviceRequests: []HostDeviceRequest{
					{
						Name:       "qat",
						DeviceName: "intel.com/qat",
					},
				},
			},
		},
	}

	assert := require.New(t)

	for i, tc := range testcases {
		v, err := parseHostDeviceInfo(tc.input)
		assert.NoError(err)
		assert.Equal(v, tc.expectation, fmt.Sprintf("failed %d: %s: expected request to match predefined objectg", i, tc.description))
	}
}
