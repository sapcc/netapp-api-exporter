package netapp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	n "github.com/pepabo/go-netapp/netapp"
)

type Client struct {
	*n.Client
	httpClient        *http.Client
	basicAuthUser     string
	basicAuthPassword string
}

func NewClient(host, username, password, version string) *Client {
	baseUrl := fmt.Sprintf("https://%s", host)
	options := &n.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}
	httpClient := &http.Client{
		Timeout: options.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !options.SSLVerify,
			},
		},
	}
	return &Client{n.NewClient(baseUrl, version, options), httpClient, username, password}
}

// Do request with internal http client. Useful to do quick checks.
func (c *Client) Do(method string, body interface{}) (*http.Response, error) {
	u, _ := c.BaseURL.Parse("/servlets/netapp.servlets.admin.XMLrequest_filer")
	buf, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "text/xml")
	}
	req.SetBasicAuth(c.basicAuthUser, c.basicAuthPassword)
	ctx, cncl := context.WithTimeout(context.Background(), c.ResponseTimeout)
	defer cncl()
	return c.httpClient.Do(req.WithContext(ctx))
}
