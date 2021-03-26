package clientbase

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"sigs.k8s.io/yaml"

	goharverrors "github.com/harvester/go-harvester/pkg/errors"
)

type APIClient struct {
	Debug      bool
	BaseURL    *url.URL
	APIVersion string
	Version    string
	PluralName string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL *url.URL, httpClient *http.Client, version string, pluralName string) *APIClient {
	return &APIClient{
		PluralName: pluralName,
		Debug:      false,
		Version:    version,
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

func (r *APIClient) BuildAPIURL() string {
	return fmt.Sprintf("%s/%s/%s", r.BaseURL, r.Version, r.PluralName)
}

func (r *APIClient) BuildResourceURL(resourceName string) string {
	if resourceName == "" {
		return r.BuildAPIURL()
	}
	return fmt.Sprintf("%s/%s", r.BuildAPIURL(), resourceName)
}

func (r *APIClient) NewRequest() *dataflow.Gout {
	return gout.New(r.HTTPClient)
}

func (r *APIClient) Create(object interface{}) (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			POST(r.BuildAPIURL()).
			SetJSON(object).
			SetHeader().
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) CreateByYAML(object interface{}) (respCode int, respBody []byte, err error) {
	var yamlData []byte
	yamlData, err = yaml.Marshal(object)
	if err != nil {
		return
	}
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			POST(r.BuildAPIURL()).
			SetBody(yamlData).
			SetCookies().
			SetHeader(gout.H{"content-type": "application/yaml"}).
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) List() (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			GET(r.BuildAPIURL()).
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) Get(resourceName string, obj ...interface{}) (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			GET(r.BuildResourceURL(resourceName)).
			SetQuery(obj...).
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) Update(resourceName string, object interface{}) (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			PUT(r.BuildResourceURL(resourceName)).
			SetJSON(object).
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) Delete(resourceName string, obj ...interface{}) (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		return r.NewRequest().
			DELETE(r.BuildResourceURL(resourceName)).
			SetQuery(obj...).
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug).
			Do()
	})
	return
}

func (r *APIClient) Action(resourceName string, action string, object interface{}) (respCode int, respBody []byte, err error) {
	err = goharverrors.RetryOnError(func() error {
		dataFlow := r.NewRequest().
			POST(fmt.Sprintf("%s?action=%s", r.BuildResourceURL(resourceName), action)).
			SetHeader().
			BindBody(&respBody).
			Code(&respCode).
			Debug(r.Debug)
		if object != nil {
			dataFlow = dataFlow.SetJSON(object)
		}
		return dataFlow.Do()
	})
	return
}
