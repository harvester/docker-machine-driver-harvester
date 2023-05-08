package harvester

import (
	"testing"
)

const (
	testCloudInitFilePath = "testdata/cloudinit_rke2.yaml"
)

func TestDriver_mergeCloudInitUserData(t *testing.T) {
	type fields struct {
		CloudConfig string
		UserData    string
	}
	tests := []struct {
		name         string
		fields       fields
		wantUserData string
		wantErr      bool
	}{
		{
			name: "empty cloud config and user data",
			fields: fields{
				CloudConfig: "",
				UserData:    "",
			},
			wantUserData: `#cloud-config
`,
			wantErr: false,
		},
		{
			name: "rke2 cloud config and empty user data",
			fields: fields{
				CloudConfig: testCloudInitFilePath,
				UserData:    "",
			},
			wantUserData: `#cloud-config
runcmd:
- sh /usr/local/custom_script/install.sh
`,
			wantErr: false,
		},
		{
			name: "rke2 cloud config and user data",
			fields: fields{
				CloudConfig: testCloudInitFilePath,
				UserData: `#cloud-config
package_update: true
packages:
- qemu-guest-agent
runcmd:
- - systemctl
- - enable
- - --now
- - qemu-guest-agent.service
`,
			},
			wantUserData: `#cloud-config
package_update: true
packages:
- qemu-guest-agent
runcmd:
- - systemctl
- - enable
- - --now
- - qemu-guest-agent.service
- sh /usr/local/custom_script/install.sh
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				CloudConfig: tt.fields.CloudConfig,
				UserData:    tt.fields.UserData,
			}
			got, _, err := d.mergeCloudInit()
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeCloudInit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantUserData {
				t.Errorf("mergeCloudInit() got = %v, want %v", got, tt.wantUserData)
			}
		})
	}
}
