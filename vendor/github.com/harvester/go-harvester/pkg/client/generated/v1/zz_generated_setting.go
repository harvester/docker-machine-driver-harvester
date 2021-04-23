package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	harv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/rancher/apiserver/pkg/types"
)

type Setting harv1.Setting

type SettingList struct {
	types.Collection
	Data []*Setting `json:"data"`
}

type SettingClient struct {
	*clientbase.APIClient
}

func newSettingClient(c *Client) *SettingClient {
	return &SettingClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "harvesterhci.io.settings"),
	}
}

func (c *SettingClient) List() (*SettingList, error) {
	var collection SettingList
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

func (c *SettingClient) Create(obj *Setting) (*Setting, error) {
	var created *Setting
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

func (c *SettingClient) Update(name string, obj *Setting) (*Setting, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Setting
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *SettingClient) Get(name string, opts ...interface{}) (*Setting, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Setting
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *SettingClient) Delete(name string, opts ...interface{}) (*Setting, error) {
	resourceName := name
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
	var obj *Setting
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
