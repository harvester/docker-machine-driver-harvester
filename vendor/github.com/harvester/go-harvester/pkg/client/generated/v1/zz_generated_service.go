package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

type Service corev1.Service

type ServiceList struct {
	types.Collection
	Data []*Service `json:"data"`
}

type ServiceClient struct {
	*clientbase.APIClient
}

func newServiceClient(c *Client) *ServiceClient {
	return &ServiceClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "services"),
	}
}

func (c *ServiceClient) List() (*ServiceList, error) {
	var collection ServiceList
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

func (c *ServiceClient) Create(obj *Service) (*Service, error) {
	var created *Service
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

func (c *ServiceClient) Update(namespace, name string, obj *Service) (*Service, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Service
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *ServiceClient) Get(namespace, name string, opts ...interface{}) (*Service, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Service
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *ServiceClient) Delete(namespace, name string, opts ...interface{}) (*Service, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Delete(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Service
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
