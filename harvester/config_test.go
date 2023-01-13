package harvester

import (
	"testing"
)

func TestCheckNetworkData(t *testing.T) {
	type args struct {
		networkDataStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				networkDataStr: "",
			},
			wantErr: false,
		},
		{
			name: "1 dhcp",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
`,
			},
			wantErr: false,
		},
		{
			name: "1 static",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: false,
		},
		{
			name: "1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "1 static without nameserver",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "2 dhcp",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: dhcp
`,
			},
			wantErr: false,
		},
		{
			name: "1 dhcp and 1 static with gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
`,
			},
			wantErr: false,
		},
		{
			name: "1 dhcp and 1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
`,
			},
			wantErr: false,
		},
		{
			name: "2 static with gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.92/24
       gateway: 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "2 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "1 static with gateway and 1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.92/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkNetworkData(tt.args.networkDataStr); (err != nil) != tt.wantErr {
				t.Errorf("CheckNetworkData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
