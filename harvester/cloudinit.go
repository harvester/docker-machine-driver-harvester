package harvester

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/harvester/harvester/pkg/builder"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	userDataHeader = `#cloud-config
`
	userDataHeaderTemplateJinja = `## template: jinja
`
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
	cloudInitNoCloudLimitSize = 2048
)

func (d *Driver) buildCloudInit() (*builder.CloudInitSource, *corev1.Secret, error) {
	cloudInitSource := &builder.CloudInitSource{
		CloudInitType: builder.CloudInitTypeNoCloud,
	}
	userData, networkData, err := d.mergeCloudInit()
	if err != nil {
		return nil, nil, err
	}
	cloudConfigSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", d.MachineName, "cloudinit"),
			Namespace: d.VMNamespace,
		},
		Data: map[string][]byte{},
	}
	if userData != "" {
		if len(userData) > cloudInitNoCloudLimitSize {
			cloudConfigSecret.Data["userdata"] = []byte(userData)
			cloudInitSource.UserDataSecretName = cloudConfigSecret.Name
		} else {
			cloudInitSource.UserData = userData
		}
	}
	if networkData != "" {
		if len(userData) > cloudInitNoCloudLimitSize {
			cloudConfigSecret.Data["networkdata"] = []byte(networkData)
			cloudInitSource.NetworkDataSecretName = cloudConfigSecret.Name
		} else {
			cloudInitSource.NetworkData = networkData
		}
	}
	if len(cloudConfigSecret.Data) == 0 {
		cloudConfigSecret = nil
	}
	return cloudInitSource, cloudConfigSecret, nil
}

func (d *Driver) mergeCloudInit() (string, string, error) {
	var (
		userData    string
		networkData string
	)
	// userData
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
	if d.CloudConfig != "" {
		cloudConfigContent, err := ioutil.ReadFile(d.CloudConfig)
		if err != nil {
			return "", "", err
		}
		userDataByte, err := mergeYaml([]byte(userData), cloudConfigContent)
		if err != nil {
			return "", "", err
		}
		userData = string(userDataByte)
	}
	if d.UserData != "" {
		userDataByte, err := mergeYaml([]byte(userData), []byte(d.UserData))
		if err != nil {
			return "", "", err
		}
		userData = string(userDataByte)
	}
	userData = userDataHeader + userData
	if strings.HasPrefix(d.UserData, userDataHeaderTemplateJinja) {
		userData = userDataHeaderTemplateJinja + userData
	}
	// networkData
	if d.NetworkData != "" {
		networkData = d.NetworkData
	}
	return userData, networkData, nil
}

func mergeYaml(dst, src []byte) ([]byte, error) {
	var (
		srcData = make(map[string]interface{})
		dstData = make(map[string]interface{})
	)
	if err := yaml.Unmarshal(src, &srcData); err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(dst, &dstData); err != nil {
		return nil, err
	}
	if err := mergo.Map(&dstData, srcData, mergo.WithAppendSlice); err != nil {
		return nil, err
	}
	return yaml.Marshal(dstData)
}
