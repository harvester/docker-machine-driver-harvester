package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	cniv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/rancher/apiserver/pkg/types"
)

type Network cniv1.NetworkAttachmentDefinition

type NetworkList struct {
	types.Collection
	Data []*Network `json:"data"`
}

type NetworkClient struct {
	*clientbase.APIClient
}

func newNetworkClient(c *Client) *NetworkClient {
	return &NetworkClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "k8s.cni.cncf.io.network-attachment-definitions"),
	}
}

func (c *NetworkClient) List() (*NetworkList, error) {
	var collection NetworkList
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

func (c *NetworkClient) Create(obj *Network) (*Network, error) {
	var created *Network
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

func (c *NetworkClient) Update(namespace, name string, obj *Network) (*Network, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Network
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *NetworkClient) Get(namespace, name string, opts ...interface{}) (*Network, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Network
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *NetworkClient) Delete(namespace, name string, opts ...interface{}) (*Network, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Delete(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Network
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
