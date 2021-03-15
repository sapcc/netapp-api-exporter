package netapp

import (
	"encoding/xml"
	"net/http"
)

type netIPSpaceRequest struct {
	Base
	CreateParams *netIPSpaceCreateParams `xml:"net-ipspaces-create,omitempty"`
	GetParams    *netIPSpaceGetParams    `xml:"net-ipspaces-get,omitempty"`
	RenameParams *netIPSpaceRenameParams `xml:"net-ipspaces-rename,omitempty"`
	DeleteParams *netIPSpaceGetParams    `xml:"net-ipspaces-destroy,omitempty"`
}

type netIPSpaceListRequest struct {
	Base
	Params struct {
		XMLName xml.Name
		NetIPSpaceOptions
	}
}

type NetIPSpaceOptions struct {
	DesiredAttributes *NetIPSpaceInfo `xml:"desired-attributes>net-ip-spaces-info,omitempty"`
	MaxRecords        int             `xml:"max-records,omitempty"`
	Query             *NetIPSpaceInfo `xml:"query>net-ipspaces-info,omitempty"`
	Tag               string          `xml:"tag,omitempty"`
}

type netIPSpaceCreateParams struct {
	IPSpace      string `xml:"ipspace"`
	ReturnRecord bool   `xml:"return-record"`
}

type netIPSpaceGetParams struct {
	IPSpace string `xml:"ipspace"`
}

type netIPSpaceRenameParams struct {
	IPSpace string `xml:"ipspace"`
	NewName string `xml:"new-name"`
}

// NetIPSpaceInfo holds newly created ipspace variables
type NetIPSpaceInfo struct {
	BroadcastDomains *[]string `xml:"broadcast-domains>broadcast-domain-name,omitempty"`
	ID               int       `xml:"id,omitempty"`
	IPSpace          string    `xml:"ipspace,omitempty"`
	Ports            *[]string `xml:"ports>net-qualified-port-name,omitempty"`
	UUID             string    `xml:"uuid,omitempty"`
	VServers         *[]string `xml:"vservers>vserver-name,omitempty"`
}

// NetIPSpaceResponse is return type for net ip space requests
type NetIPSpaceResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		NetIPSpaceInfo       `xml:",innerxml"`
		NetIPSpaceCreateInfo *NetIPSpaceInfo `xml:"result>net-ipspaces-info"`
	} `xml:"results"`
}

type NetIPSpaceListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		Info       []NetIPSpaceInfo `xml:"attributes-list>net-ipspaces-info"`
		NumRecords string           `xml:"num-records"`
	} `xml:"results"`
}

// CreateIPSpace creates a new ipspace on the cluster
func (n Net) CreateIPSpace(name string, returnRecord bool) (*NetIPSpaceResponse, *http.Response, error) {
	req := n.newNetIPSpaceRequest()
	req.CreateParams = &netIPSpaceCreateParams{
		IPSpace:      name,
		ReturnRecord: returnRecord,
	}
	return n.newNetIPSpaceResponse(req)
}

// GetIPSpace grabs data for an ip space
func (n Net) GetIPSpace(name string) (*NetIPSpaceResponse, *http.Response, error) {
	req := n.newNetIPSpaceRequest()
	req.GetParams = &netIPSpaceGetParams{
		IPSpace: name,
	}

	return n.newNetIPSpaceResponse(req)
}

func (n Net) ListIPSpaces(query *NetIPSpaceInfo) (*NetIPSpaceListResponse, *http.Response, error) {
	req := &netIPSpaceListRequest{
		Base: n.Base,
	}
	req.Params.XMLName = xml.Name{Local: "net-ipspaces-get-iter"}
	req.Params.MaxRecords = 20
	req.Params.NetIPSpaceOptions.Query = query

	r := NetIPSpaceListResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// RenameIPSpace changes the name of an ipspace
func (n Net) RenameIPSpace(name string, newName string) (*NetIPSpaceResponse, *http.Response, error) {
	req := n.newNetIPSpaceRequest()
	req.RenameParams = &netIPSpaceRenameParams{
		IPSpace: name,
		NewName: newName,
	}

	return n.newNetIPSpaceResponse(req)
}

// DeleteIPSpace deletes an IPSpace
func (n Net) DeleteIPSpace(name string) (*NetIPSpaceResponse, *http.Response, error) {
	req := n.newNetIPSpaceRequest()
	req.DeleteParams = &netIPSpaceGetParams{
		IPSpace: name,
	}

	return n.newNetIPSpaceResponse(req)
}

func (n Net) newNetIPSpaceRequest() *netIPSpaceRequest {
	return &netIPSpaceRequest{
		Base: n.Base,
	}
}

func (n Net) newNetIPSpaceResponse(req *netIPSpaceRequest) (*NetIPSpaceResponse, *http.Response, error) {
	r := NetIPSpaceResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}
