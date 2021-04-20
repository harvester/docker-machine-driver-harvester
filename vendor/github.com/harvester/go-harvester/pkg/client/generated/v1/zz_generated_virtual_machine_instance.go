package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

type VirtualMachineInstance kubevirtv1.VirtualMachineInstance

type VirtualMachineInstanceList struct {
	types.Collection
	Data []*VirtualMachineInstance `json:"data"`
}

type VirtualMachineInstanceClient struct {
	*clientbase.APIClient
}

func newVirtualMachineInstanceClient(c *Client) *VirtualMachineInstanceClient {
	return &VirtualMachineInstanceClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "kubevirt.io.virtualmachineinstance"),
	}
}

func (c *VirtualMachineInstanceClient) List() (*VirtualMachineInstanceList, error) {
	var collection VirtualMachineInstanceList
	respCode, respBody, err := c.APIClient.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (c *VirtualMachineInstanceClient) Create(obj *VirtualMachineInstance) (*VirtualMachineInstance, error) {
	var created *VirtualMachineInstance
	respCode, respBody, err := c.APIClient.Create(obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusCreated {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &created)
	return created, nil
}

func (c *VirtualMachineInstanceClient) Update(namespace, name string, obj *VirtualMachineInstance) (*VirtualMachineInstance, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *VirtualMachineInstance
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *VirtualMachineInstanceClient) Get(namespace, name string, opts ...interface{}) (*VirtualMachineInstance, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *VirtualMachineInstance
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *VirtualMachineInstanceClient) Delete(namespace, name string, opts ...interface{}) (*VirtualMachineInstance, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Delete(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode == http.StatusNoContent {
		return nil, nil
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *VirtualMachineInstance
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
