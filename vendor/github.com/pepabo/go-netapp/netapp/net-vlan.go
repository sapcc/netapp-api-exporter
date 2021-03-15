package netapp

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type netVlanRequest struct {
	Base
	Params struct {
		XMLName         xml.Name
		NetVlanInfo     `xml:",innerxml"`
		*NetVlanOptions `xml:",innerxml"`

		VlanInfo *NetVlanInfo `xml:"vlan-info,omitempty"`
	}
}

// NetVlanOptions get/list options for getting vlans
type NetVlanOptions struct {
	DesiredAttributes *NetVlanInfo `xml:"desired-attributes>vlan-info,omitempty"`
	MaxRecords        int          `xml:"max-records,omitempty"`
	Query             *NetVlanInfo `xml:"query>vlan-info,omitempty"`
	Tag               string       `xml:"tag,omitempty"`
}

// NetVlanInfo is the Vlan data
type NetVlanInfo struct {
	InterfaceName   string `xml:"interface-name,omitempty"`
	Node            string `xml:"node,omitempty"`
	ParentInterface string `xml:"parent-interface,omitempty"`
	VlanID          int    `xml:"vlanid,omitempty"`
}

// NetVlanResponse returns results for a single vlan
type NetVlanResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		Info NetVlanInfo `xml:"attributes>vlan-info"`
	} `xml:"results"`
}

// NetVlanListResponse returns results for a list of vlans
type NetVlanListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		Info []NetVlanInfo `xml:"attributes-list>vlan-info"`
	} `xml:"results"`
}

// ToString converts to string, ie test-cluster-01-01:a0a-3555
func (v *NetVlanInfo) ToString() string {
	return fmt.Sprintf("%s:%s-%d", v.Node, v.ParentInterface, v.VlanID)
}

// CreateVlan creates a new vlan
func (n Net) CreateVlan(options *NetVlanInfo) (*SingleResultResponse, *http.Response, error) {
	req := n.newNetVlanRequest()
	req.Params.XMLName = xml.Name{Local: "net-vlan-create"}

	options.InterfaceName = ""
	req.Params.VlanInfo = options

	r := SingleResultResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// GetVlan grabs a single named broadcast domain
func (n Net) GetVlan(interfaceName string, node string) (*NetVlanResponse, *http.Response, error) {
	req := n.newNetVlanRequest()
	req.Params.XMLName = xml.Name{Local: "net-vlan-get"}
	req.Params.InterfaceName = interfaceName
	req.Params.Node = node
	r := NetVlanResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// ListVlans lists all vlans that match info query
func (n Net) ListVlans(info *NetVlanInfo) (*NetVlanListResponse, *http.Response, error) {
	req := n.newNetVlanRequest()

	req.Params.XMLName = xml.Name{Local: "net-vlan-get-iter"}
	req.Params.NetVlanOptions = &NetVlanOptions{
		Query:      info,
		MaxRecords: 20,
	}

	r := NetVlanListResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// DeleteVlan removes vlan from existence
func (n Net) DeleteVlan(options *NetVlanInfo) (*SingleResultResponse, *http.Response, error) {
	req := n.newNetVlanRequest()
	req.Params.XMLName = xml.Name{Local: "net-vlan-delete"}
	options.InterfaceName = ""
	req.Params.VlanInfo = options
	r := SingleResultResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

func (n Net) newNetVlanRequest() *netVlanRequest {
	return &netVlanRequest{
		Base: n.Base,
	}
}
