package harvester

import "testing"

func TestFormatLabelValue(t *testing.T) {
	testCases := []struct {
		desc       string
		labelValue string
		expected   string
	}{
		{
			desc:       "return empty string if label value is empty",
			labelValue: "",
			expected:   "",
		},
		{
			desc:       "return label value unchanged if it's less than 63 characters",
			labelValue: "machineSetName",
			expected:   "machineSetName",
		},
		{
			desc:       "return hashed label value if more than 63 characters",
			labelValue: "machineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetNamemachineSetName",
			expected:   "hash_FR_ghQ_z",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(*testing.T) {
			actual, err := formatLabelValue(tc.labelValue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.expected != actual {
				t.Errorf("test case failed: %s. expected %s, got %s", tc.desc, tc.expected, actual)
			}
		})
	}
}
