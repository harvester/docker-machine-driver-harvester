package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

type Node corev1.Node

type NodeList struct {
	types.Collection
	Data []*Node `json:"data"`
}

type NodeClient struct {
	*clientbase.APIClient
}

func newNodeClient(c *Client) *NodeClient {
	return &NodeClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "nodes"),
	}
}

func (c *NodeClient) List() (*NodeList, error) {
	var collection NodeList
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

func (c *NodeClient) Create(obj *Node) (*Node, error) {
	var created *Node
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

func (c *NodeClient) Update(name string, obj *Node) (*Node, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Node
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *NodeClient) Get(name string, opts ...interface{}) (*Node, error) {
	resourceName := name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Node
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *NodeClient) Delete(name string, opts ...interface{}) (*Node, error) {
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
	var obj *Node
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
