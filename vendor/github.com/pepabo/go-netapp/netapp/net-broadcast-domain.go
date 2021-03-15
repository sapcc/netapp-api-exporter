package netapp

import (
	"encoding/xml"
	"net/http"
)

type netBroadcastDomainRequest struct {
	Base
	Params struct {
		XMLName xml.Name
		NetBroadcastDomainOptions
		NetBroadcastDomainCreateOptions `xml:",innerxml"`
	}
}

// NetBroadcastDomainOptions get/list options for getting broadcast domains
type NetBroadcastDomainOptions struct {
	DesiredAttributes *NetBroadcastDomainInfo `xml:"desired-attributes,omitempty"`
	MaxRecords        int                     `xml:"max-records,omitempty"`
	Query             *NetBroadcastDomainInfo `xml:"query,omitempty"`
	Tag               string                  `xml:"tag,omitempty"`
}

// NetBroadcastDomainInfo is the Broadcast Domain data
type NetBroadcastDomainInfo struct {
	BroadcastDomain          string               `xml:"broadcast-domain,omitempty"`
	FailoverGroups           []string             `xml:"failover-groups>failover-group"`
	IPSpace                  string               `xml:"ipspace,omitempty"`
	MTU                      int                  `xml:"mtu,omitempty"`
	CombinedPortUpdateStatus string               `xml:"port-update-status-combined,omitempty"`
	Ports                    *[]NetPortUpdateInfo `xml:"ports>port-info"`
	SubnetNames              []string             `xml:"subnet-names>subnet-name"`
}

// NetBroadcastDomainCreateOptions used for creating new Broadcast Domain
type NetBroadcastDomainCreateOptions struct {
	BroadcastDomain string    `xml:"broadcast-domain"`
	IPSpace         string    `xml:"ipspace"`
	MTU             int       `xml:"mtu,omitempty"`
	Ports           *[]string `xml:"ports>net-qualified-port-name,omitempty"`
}

// NetPortUpdateInfo is port info for the broadcast domain
type NetPortUpdateInfo struct {
	Port                    string `xml:"port"`
	PortUpdateStatus        string `xml:"port-update-status"`
	PortUpdateStatusDetails string `xml:"port-update-status-details"`
}

// NetBroadcastDomainResponse returns results for broadcast domains
type NetBroadcastDomainResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		Info NetBroadcastDomainInfo `xml:"attributes>net-port-broadcast-domain-info"`
	} `xml:"results"`
}

// NetBroadcastDomainCreateResponse returns result of creating a new broadcast domain
type NetBroadcastDomainCreateResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		CombinedPortUpdateStatus string `xml:"port-update-status-combined"`
	} `xml:"results"`
}

// CreateBroadcastDomain creates a new broadcast domain
func (n Net) CreateBroadcastDomain(createOptions *NetBroadcastDomainCreateOptions) (*NetBroadcastDomainCreateResponse, *http.Response, error) {
	req := n.newNetBroadcastDomainRequest()
	req.Params.XMLName = xml.Name{Local: "net-port-broadcast-domain-create"}
	req.Params.NetBroadcastDomainCreateOptions = *createOptions

	r := NetBroadcastDomainCreateResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// GetBroadcastDomain grabs a single named broadcast domain
func (n Net) GetBroadcastDomain(domain string, ipSpace string) (*NetBroadcastDomainResponse, *http.Response, error) {
	req := n.newNetBroadcastDomainRequest()
	req.Params.XMLName = xml.Name{Local: "net-port-broadcast-domain-get"}
	req.Params.BroadcastDomain = domain
	req.Params.IPSpace = ipSpace
	r := NetBroadcastDomainResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

func (n Net) DeleteBroadcastDomain(domain string, ipSpace string) (*NetBroadcastDomainCreateResponse, *http.Response, error) {
	req := n.newNetBroadcastDomainRequest()
	req.Params.XMLName = xml.Name{Local: "net-port-broadcast-domain-destroy"}
	req.Params.BroadcastDomain = domain
	req.Params.IPSpace = ipSpace

	r := NetBroadcastDomainCreateResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

func (n Net) newNetBroadcastDomainRequest() *netBroadcastDomainRequest {
	return &netBroadcastDomainRequest{
		Base: n.Base,
	}
}
