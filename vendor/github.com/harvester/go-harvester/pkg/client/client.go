package client

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	v1client "github.com/harvester/go-harvester/pkg/client/generated/v1"
)

type Client struct {
	*v1client.Client

	Auth *AuthClient
}

func New(harvesterURL string, transport *http.Transport) (*Client, error) {
	baseURL, err := url.Parse(harvesterURL)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	if transport == nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	httpClient := &http.Client{
		Jar:       jar,
		Transport: transport,
	}
	c := &Client{
		Client: v1client.New(baseURL, httpClient),
		Auth:   newAuthClient(baseURL, httpClient),
	}

	return c, nil
}
