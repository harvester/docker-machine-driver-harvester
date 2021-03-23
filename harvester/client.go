package harvester

import (
	"fmt"

	goharv "github.com/harvester/go-harvester/pkg/client"
)

func (d *Driver) getClient() (*goharv.Client, error) {
	if d.client == nil {
		c, err := d.login()
		if err != nil {
			return nil, err
		}
		d.client = c
	}
	return d.client, nil
}

func (d *Driver) login() (*goharv.Client, error) {
	c, err := goharv.New(fmt.Sprintf("https://%s:%d", d.Host, d.Port), nil)
	if err != nil {
		return nil, err
	}
	if err = c.Auth.Login(d.Username, d.Password); err != nil {
		return nil, err
	}

	return c, nil
}
