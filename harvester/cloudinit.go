package harvester

import "fmt"

const (
	userDataHeader            = `#cloud-config`
	userDataAddQemuGuestAgent = `
package_update: true
packages:
- qemu-guest-agent
runcmd:
- [systemctl, enable, --now, qemu-guest-agent]`
	userDataPasswordTemplate = `
user: %s
password: %s
chpasswd: { expire: False }
ssh_pwauth: True`

	userDataSSHKeyTemplate = `
ssh_authorized_keys:
- >-
  %s`
	userDataAddDockerGroupSSHKeyTemplate = `
groups:
- docker
users:
- name: %s
  sudo: ALL=(ALL) NOPASSWD:ALL
  groups: sudo, docker
  shell: /bin/bash
  ssh_authorized_keys:
  - >-
    %s`
)

func (d *Driver) createCloudInit() (userData string, networkData string) {
	// userData
	userData = userDataHeader
	if d.NetworkType != networkTypePod {
		// need qemu guest agent to get ip
		userData += userDataAddQemuGuestAgent
	}
	if d.SSHPassword != "" {
		userData += fmt.Sprintf(userDataPasswordTemplate, d.SSHUser, d.SSHPassword)
	}
	if d.SSHPublicKey != "" {
		if d.AddUserToDockerGroup {
			userData += fmt.Sprintf(userDataAddDockerGroupSSHKeyTemplate, d.SSHUser, d.SSHPublicKey)
		} else {
			userData += fmt.Sprintf(userDataSSHKeyTemplate, d.SSHPublicKey)
		}
	}
	return
}
