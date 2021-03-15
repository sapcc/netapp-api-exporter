package netapp

import (
	"encoding/xml"
	"net/http"
)

type ClusterIdentity struct {
	Base
	Params struct {
		XMLName xml.Name
		*ClusterIdentityOptions
	}
}

type ClusterIdentityInfo struct {
	ClusterContact      string `xml:"cluster-contact,omitempty"`
	ClusterLocation     string `xml:"cluster-location"`
	ClusterName         string `xml:"cluster-name"`
	ClusterSerialNumber string `xml:"cluster-serial-number"`
	RdbUuid             string `xml:"rdb-uuid"`
	UUID                string `xml:"uuid"`
}

type ClusterIdentityOptions struct {
	DesiredAttributes *ClusterIdentityInfo `xml:"desired-attributes,omitempty"`
}

type ClusterIdentityResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		ClusterIdentityInfo []ClusterIdentityInfo `xml:"attributes>cluster-identity-info"`
	} `xml:"results"`
}

func (c *ClusterIdentity) List(options *ClusterIdentityOptions) (*ClusterIdentityResponse, *http.Response, error) {
	c.Params.XMLName = xml.Name{Local: "cluster-identity-get"}
	c.Params.ClusterIdentityOptions = options
	r := ClusterIdentityResponse{}
	res, err := c.get(c, &r)
	return &r, res, err
}
