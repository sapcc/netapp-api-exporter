package netapp

import (
	"fmt"
	"time"

	n "github.com/pepabo/go-netapp/netapp"
)

type Client struct {
	*n.Client
}

func NewClient(host, username, password, version string) *Client {
	_url := "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	url := fmt.Sprintf(_url, host)

	opts := &n.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}

	return &Client{n.NewClient(url, version, opts)}
}
