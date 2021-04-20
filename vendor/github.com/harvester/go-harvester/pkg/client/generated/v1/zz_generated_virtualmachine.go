package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

type VirtualMachine kubevirtv1.VirtualMachine

type VirtualMachineList struct {
	types.Collection
	Data []*VirtualMachine `json:"data"`
}

type VirtualMachineClient struct {
	*clientbase.APIClient
}

func newVirtualMachineClient(c *Client) *VirtualMachineClient {
	return &VirtualMachineClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "kubevirt.io.virtualmachines"),
	}
}

func (c *VirtualMachineClient) List() (*VirtualMachineList, error) {
	var collection VirtualMachineList
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

func (c *VirtualMachineClient) Create(obj *VirtualMachine) (*VirtualMachine, error) {
	var created *VirtualMachine
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

func (c *VirtualMachineClient) Update(namespace, name string, obj *VirtualMachine) (*VirtualMachine, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *VirtualMachine
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *VirtualMachineClient) Get(namespace, name string, opts ...interface{}) (*VirtualMachine, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *VirtualMachine
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *VirtualMachineClient) Delete(namespace, name string, opts ...interface{}) (*VirtualMachine, error) {
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
	var obj *VirtualMachine
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *VirtualMachineClient) AbortMigration(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "abortMigration", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Backup(namespace, name string, backupInput interface{}) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "backup", backupInput)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) EjectCdRom(namespace, name string, ejectCdRomActionInput interface{}) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "ejectCdRom", ejectCdRomActionInput)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Migrate(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "migrate", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Pause(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "pause", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Restart(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "restart", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Restore(namespace, name string, restoreInput interface{}) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "restore", restoreInput)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Start(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "start", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Stop(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "stop", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}

func (c *VirtualMachineClient) Unpause(namespace, name string) error {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Action(resourceName, "unpause", nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return errors.NewResponseError(respCode, respBody)
	}
	return nil
}
