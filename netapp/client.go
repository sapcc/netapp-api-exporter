package netapp

import (
	"fmt"
	"time"

	n "github.com/pepabo/go-netapp/netapp"
)

type Filer struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	AvailabilityZone string `yaml:"availability_zone"`
	Version          string `yaml:"version"`
}

type Client struct {
	*n.Client
}

func NewClient(f *Filer) *Client {
	return &Client{
		newNetappClient(f.Host, f.Username, f.Password, f.Version),
	}
}

func newNetappClient(host, username, password, version string) *n.Client {
	_url := "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	url := fmt.Sprintf(_url, host)

	opts := &n.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}

	return n.NewClient(url, version, opts)
}
