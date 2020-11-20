package netapp

import (
	"encoding/xml"
)

func (c *Client) CheckCluster() (statusCode int, err error) {
	body := c.ClusterIdentity
	body.Params.XMLName = xml.Name{Local: "cluster-identity-get"}
	resp, err := c.Do("POST", body)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return
}
