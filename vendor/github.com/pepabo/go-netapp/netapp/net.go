package netapp

import (
	"encoding/xml"
	"net/http"
)

type Net struct {
	Base
}

type NetPortQuery struct {
	NetPortInfo *NetPortInfo `xml:"net-port-info,omitempty"`
}

type NetPortOptions struct {
	DesiredAttributes *NetPortQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int           `xml:"max-records,omitempty"`
	Query             *NetPortQuery `xml:"query,omitempty"`
	Tag               string        `xml:"tag,omitempty"`
}

type NetPortInfo struct {
	AdministrativeDuplex          string `xml:"administrative-duplex,omitempty"`
	AdministrativeFlowcontrol     string `xml:"administrative-flowcontrol,omitempty"`
	AdministrativeSpeed           string `xml:"administrative-speed,omitempty"`
	AutorevertDelay               int    `xml:"autorevert-delay,omitempty"`
	IfgrpDistributionFunction     string `xml:"ifgrp-distribution-function,omitempty"`
	IfgrpMode                     string `xml:"ifgrp-mode,omitempty"`
	IfgrpNode                     string `xml:"ifgrp-node,omitempty"`
	IfgrpPort                     string `xml:"ifgrp-port,omitempty"`
	IsAdministrativeAutoNegotiate bool   `xml:"is-administrative-auto-negotiate,omitempty"`
	IsAdministrativeUp            bool   `xml:"is-administrative-up,omitempty"`
	IsOperationalAutoNegotiate    bool   `xml:"is-operational-auto-negotiate,omitempty"`
	LinkStatus                    string `xml:"link-status,omitempty"`
	MacAddress                    string `xml:"mac-address,omitempty"`
	Mtu                           int    `xml:"mtu,omitempty"`
	Node                          string `xml:"node,omitempty"`
	OperationalDuplex             string `xml:"operational-duplex,omitempty"`
	OperationalFlowcontrol        string `xml:"operational-flowcontrol,omitempty"`
	OperationalSpeed              string `xml:"operational-speed,omitempty"`
	Port                          string `xml:"port,omitempty"`
	PortType                      string `xml:"port-type,omitempty"`
	RemoteDeviceId                string `xml:"remote-device-id,omitempty"`
	Role                          string `xml:"role,omitempty"`
	VlanId                        int    `xml:"vlan-id,omitempty"`
	VlanNode                      string `xml:"vlan-node,omitempty"`
	VlanPort                      string `xml:"vlan-port,omitempty"`
}

type NetPortGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			NetPortAttributes []NetPortInfo `xml:"net-port-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type NetPortPageResponse struct {
	Response    *NetPortGetIterResponse
	Error       error
	RawResponse *http.Response
}

type NetPortPageHandler func(NetPortPageResponse) (shouldContinue bool)

func (n *Net) NetPortGetIter(options *NetPortOptions) (*NetPortGetIterResponse, *http.Response, error) {
	params := newNetPortGetIterParams(options, n.Base)
	r := NetPortGetIterResponse{}
	res, err := n.get(params, &r)
	return &r, res, err
}

func (n *Net) NetPortGetAll(options *NetPortOptions, fn NetPortPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		netPortGetIterResponse, res, err := n.NetPortGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(NetPortPageResponse{Response: netPortGetIterResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = netPortGetIterResponse.Results.NextTag
			requestOptions = &NetPortOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}

type netPortGetIterParams struct {
	Base
	Params struct {
		XMLName xml.Name
		*NetPortOptions
	}
}

func newNetPortGetIterParams(options *NetPortOptions, base Base) *netPortGetIterParams {
	params := netPortGetIterParams{
		Base: base,
	}
	params.Params.XMLName = xml.Name{Local: "net-port-get-iter"}
	params.Params.NetPortOptions = options
	return &params
}

type NetInterfaceQuery struct {
	NetInterfaceInfo *NetInterfaceInfo `xml:"net-interface-info,omitempty"`
}

type NetInterfaceOptions struct {
	DesiredAttributes *NetInterfaceQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                `xml:"max-records,omitempty"`
	Query             *NetInterfaceQuery `xml:"query,omitempty"`
	Tag               string             `xml:"tag,omitempty"`
}

type NetInterfaceInfo struct {
	Address              string    `xml:"address,omitempty"`
	AdministrativeStatus string    `xml:"administrative-status,omitempty"`
	Comment              string    `xml:"comment,omitempty"`
	DataProtocols        *[]string `xml:"data-protocols>data-protocol"`
	CurrentNode          string    `xml:"current-node,omitempty"`
	CurrentPort          string    `xml:"current-port,omitempty"`
	DnsDomainName        string    `xml:"dns-domain-name,omitempty"`
	FailoverGroup        string    `xml:"failover-group,omitempty"`
	FailoverPolicy       string    `xml:"failover-policy,omitempty"`
	FirewallPolicy       string    `xml:"firewall-policy,omitempty"`
	HomeNode             string    `xml:"home-node,omitempty"`
	HomePort             string    `xml:"home-port,omitempty"`
	InterfaceName        string    `xml:"interface-name,omitempty"`
	IsAutoRevert         bool      `xml:"is-auto-revert,omitempty"`
	IsHome               bool      `xml:"is-home,omitempty"`
	IsIpv4LinkLocal      bool      `xml:"is-ipv4-link-local,omitempty"`
	Netmask              string    `xml:"netmask,omitempty"`
	NetmaskLength        int       `xml:"netmask-length,omitempty"`
	OperationalStatus    string    `xml:"operational-status,omitempty"`
	Role                 string    `xml:"role,omitempty"`
	RoutingGroupName     string    `xml:"routing-group-name,omitempty"`
	UseFailoverGroup     string    `xml:"use-failover-group,omitempty"`
	Vserver              string    `xml:"vserver,omitempty"`
}

type NetInterfaceGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			NetInterfaceAttributes []NetInterfaceInfo `xml:"net-interface-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type NetInterfacePageResponse struct {
	Response    *NetInterfaceGetIterResponse
	Error       error
	RawResponse *http.Response
}

type NetInterfacePageHandler func(NetInterfacePageResponse) (shouldContinue bool)

// CreateNetInterface creates a new network interface
func (n Net) CreateNetInterface(options *NetInterfaceInfo) (*SingleResultResponse, *http.Response, error) {
	req := netInterfaceCreateRequest{
		Base: n.Base,
	}
	req.Params.XMLName = xml.Name{Local: "net-interface-create"}
	req.Params.NetInterfaceInfo = *options
	r := SingleResultResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// DeleteNetInterface removes a LIF from a vserver
func (n Net) DeleteNetInterface(vServerName string, lif string) (*SingleResultResponse, *http.Response, error) {
	req := netInterfaceCreateRequest{
		Base: n.Base,
	}
	req.Params.XMLName = xml.Name{Local: "net-interface-delete"}
	req.Params.Vserver = vServerName
	req.Params.InterfaceName = lif
	r := SingleResultResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

func (n *Net) NetInterfaceGetIter(options *NetInterfaceOptions) (*NetInterfaceGetIterResponse, *http.Response, error) {
	params := newNetInterfaceGetIterParams(options, n.Base)
	r := NetInterfaceGetIterResponse{}
	res, err := n.get(params, &r)
	return &r, res, err
}

func (n *Net) NetInterfaceGetAll(options *NetInterfaceOptions, fn NetInterfacePageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		netInterfaceGetIterResponse, res, err := n.NetInterfaceGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(NetInterfacePageResponse{Response: netInterfaceGetIterResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = netInterfaceGetIterResponse.Results.NextTag
			requestOptions = &NetInterfaceOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}

type netInterfaceCreateRequest struct {
	Base
	Params struct {
		XMLName          xml.Name
		NetInterfaceInfo `xml:",innerxml"`
	}
}

type netInterfaceGetIterParams struct {
	Base
	Params struct {
		XMLName xml.Name
		*NetInterfaceOptions
	}
}

func newNetInterfaceGetIterParams(options *NetInterfaceOptions, base Base) *netInterfaceGetIterParams {
	params := netInterfaceGetIterParams{
		Base: base,
	}
	params.Params.XMLName = xml.Name{Local: "net-interface-get-iter"}
	params.Params.NetInterfaceOptions = options
	return &params
}

// NetRoutingGroupRouteInfo holds route information
type NetRoutesInfo struct {
	AddressFamily      string `xml:"address-family,omitempty"`
	DestinationAddress string `xml:"destination"`
	GatewayAddress     string `xml:"gateway"`
	Metric             int    `xml:"metric,omitempty"`
	ReturnRecord       bool   `xml:"return-record,omitempty"`
	VServer            string `xml:"vserver,omitempty"`
}

type netRoutesRequest struct {
	Base
	Params struct {
		XMLName       xml.Name
		NetRoutesInfo `xml:",innerxml"`
	}
}

type NetRoutesResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		Info NetRoutesInfo `xml:"result>net-vs-routes-info"`
	} `xml:"results"`
}

// CreateRoute creates a new route for Routing Group
func (n Net) CreateRoute(vServerName string, options *NetRoutesInfo) (*NetRoutesResponse, *http.Response, error) {
	req := netRoutesRequest{
		Base: n.Base,
	}
	req.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "net-routes-create"}
	req.Params.NetRoutesInfo = *options
	r := NetRoutesResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}

// DeleteRoute creates a new route for Routing Group
func (n Net) DeleteRoute(vServerName string, destination string, gateway string) (*SingleResultResponse, *http.Response, error) {
	req := netRoutesRequest{
		Base: n.Base,
	}
	req.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "net-routes-destroy"}
	req.Params.DestinationAddress = destination
	req.Params.GatewayAddress = gateway
	r := SingleResultResponse{}
	res, err := n.get(req, &r)
	return &r, res, err
}
