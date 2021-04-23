package client

import (
	"encoding/json"
	"net/http"

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
	harv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/rancher/apiserver/pkg/types"
)

type Image harv1.VirtualMachineImage

type ImageList struct {
	types.Collection
	Data []*Image `json:"data"`
}

type ImageClient struct {
	*clientbase.APIClient
}

func newImageClient(c *Client) *ImageClient {
	return &ImageClient{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "v1", "harvesterhci.io.virtualmachineimages"),
	}
}

func (c *ImageClient) List() (*ImageList, error) {
	var collection ImageList
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

func (c *ImageClient) Create(obj *Image) (*Image, error) {
	var created *Image
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

func (c *ImageClient) Update(namespace, name string, obj *Image) (*Image, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *Image
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (c *ImageClient) Get(namespace, name string, opts ...interface{}) (*Image, error) {
	resourceName := namespace + "/" + name
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *Image
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

func (c *ImageClient) Delete(namespace, name string, opts ...interface{}) (*Image, error) {
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
	var obj *Image
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
