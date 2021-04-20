package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvesterhci.io/v1beta1"
)

type Keypair harv1.KeyPair

type KeypairList struct {
	types.Collection
	Data []*Keypair `json:"data"`
}

type KeypairClient struct {
	*clientbase.APIClient
}

func newKeypairClient(c *Client) *KeypairClient {
	return &KeypairClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "harvesterhci.io.keypairs"),
	}
}

func (c *KeypairClient) List() (*KeypairList, error) {
	var collection KeypairList
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

func (c *KeypairClient) Create(obj *Keypair) (*Keypair, error) {
	var created *Keypair
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

func (c *KeypairClient) Update(namespace, name string, obj *Keypair) (*Keypair, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Keypair
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *KeypairClient) Get(namespace, name string, opts ...interface{}) (*Keypair, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Keypair
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *KeypairClient) Delete(namespace, name string, opts ...interface{}) (*Keypair, error) {
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
	var obj *Keypair
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
