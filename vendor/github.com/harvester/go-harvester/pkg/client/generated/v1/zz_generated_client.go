package client

import (
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL

	Nodes                   *NodeClient
	Services                *ServiceClient
	Users                   *UserClient
	VirtualMachines         *VirtualMachineClient
	Images                  *ImageClient
	Keypairs                *KeypairClient
	Settings                *SettingClient
	Networks                *NetworkClient
	VirtualMachineInstances *VirtualMachineInstanceClient
	Volumes                 *VolumeClient
}

func New(baseURL *url.URL, httpClient *http.Client) *Client {

	c := &Client{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}

	c.Nodes = newNodeClient(c)
	c.Services = newServiceClient(c)
	c.Users = newUserClient(c)
	c.VirtualMachines = newVirtualMachineClient(c)
	c.Images = newImageClient(c)
	c.Keypairs = newKeypairClient(c)
	c.Settings = newSettingClient(c)
	c.Networks = newNetworkClient(c)
	c.VirtualMachineInstances = newVirtualMachineInstanceClient(c)
	c.Volumes = newVolumeClient(c)

	return c
}
