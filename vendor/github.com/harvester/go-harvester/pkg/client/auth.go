package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/guonaihong/gout"
	"github.com/rancher/wrangler/pkg/slice"

	"github.com/harvester/go-harvester/pkg/clientbase"
	goharverrors "github.com/harvester/go-harvester/pkg/errors"
)

func UnmarshalAuthModes(data []byte) (AuthModes, error) {
	var r AuthModes
	err := json.Unmarshal(data, &r)
	return r, err
}

type AuthModes struct {
	Modes []string `json:"modes"`
}

type AuthClient struct {
	v1AuthMode       *clientbase.APIClient
	v1Auth           *clientbase.APIClient
	v3localProviders *clientbase.APIClient
}

func newAuthClient(baseURL *url.URL, httpClient *http.Client) *AuthClient {
	return &AuthClient{
		v1AuthMode:       clientbase.NewAPIClient(baseURL, httpClient, "v1-public", "auth-modes"),
		v1Auth:           clientbase.NewAPIClient(baseURL, httpClient, "v1-public", "auth"),
		v3localProviders: clientbase.NewAPIClient(baseURL, httpClient, "v3-public", "localProviders"),
	}
}

func (c *AuthClient) Login(username, password string) error {
	respCode, respBody, err := c.v1AuthMode.List()
	if err != nil {
		return err
	}
	if respCode != http.StatusOK {
		return goharverrors.NewResponseError(respCode, respBody)
	}
	authModes, err := UnmarshalAuthModes(respBody)
	if err != nil {
		return err
	}
	if len(authModes.Modes) == 1 && authModes.Modes[0] == "rancher" {
		respCode, respBody, err = c.v3localProviders.Action("local", "login", gout.H{
			"username":     username,
			"password":     password,
			"ttl":          57600000,
			"description":  "UI Session",
			"responseType": "cookie",
		})
	} else if slice.ContainsString(authModes.Modes, "localUser") {
		respCode, respBody, err = c.v1Auth.Action("", "login", gout.H{
			"username": username,
			"password": password,
		})
	} else {
		err = errors.New("authMode not localUser or rancher")
	}
	if err != nil {
		return err
	}
	if respCode != http.StatusOK {
		return goharverrors.NewResponseError(respCode, respBody)
	}
	return nil
}
