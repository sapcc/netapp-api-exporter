package netapp

import (
	"fmt"

	n "github.com/pepabo/go-netapp/netapp"
)

func (c *Client) GetSystemVersion() (string, error) {
	opts := &n.NodeDetailOptions{}
	resp, httpResp, err := c.System.List(opts)
	if err != nil {
		return "", err
	}
	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("http request failed with %v", httpResp.Status)
	}
	if len(resp.Results.NodeDetails) == 0 {
		return "", fmt.Errorf("failed to get node details")
	}
	return resp.Results.NodeDetails[0].ProductVersion, nil
}
