package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvesterhci.io/v1beta1"
)

type User harv1.User

type UserList struct {
	types.Collection
	Data []*User `json:"data"`
}

type UserClient struct {
	*clientbase.APIClient
}

func newUserClient(c *Client) *UserClient {
	return &UserClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "harvesterhci.io.users"),
	}
}

func (c *UserClient) List() (*UserList, error) {
	var collection UserList
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

func (c *UserClient) Create(obj *User) (*User, error) {
	var created *User
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

func (c *UserClient) Update(name string, obj *User) (*User, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *User
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *UserClient) Get(name string, opts ...interface{}) (*User, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *User
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *UserClient) Delete(name string, opts ...interface{}) (*User, error) {
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
	var obj *User
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
