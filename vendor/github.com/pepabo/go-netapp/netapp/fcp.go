package netapp

import (
	"encoding/xml"
	"net/http"
)

type Fcp struct {
	Base
	Params struct {
		XMLName xml.Name
		*FcpAdapterConfigOptions
	}
}

type FcpAdapterConfigQuery struct {
	FcpAdapterConfigInfo *FcpAdapterConfigInfo `xml:"net-port-info,omitempty"`
}

type FcpAdapterConfigOptions struct {
	DesiredAttributes *FcpAdapterConfigQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                    `xml:"max-records,omitempty"`
	Query             *FcpAdapterConfigQuery `xml:"query,omitempty"`
	Tag               string                 `xml:"tag,omitempty"`
}

type FcpAdapterConfigInfo struct {
	Adapter               string `xml:"adapter"`
	CacheLineSize         int    `xml:"cache-line-size"`
	ConnectionEstablished string `xml:"connection-established"`
	DataLinkRate          int    `xml:"data-link-rate"`
	ExternalGbicEnabled   bool   `xml:"external-gbic-enabled"`
	FabricEstablished     bool   `xml:"fabric-established"`
	FirmwareRev           string `xml:"firmware-rev"`
	HardwareRev           string `xml:"hardware-rev"`
	InfoName              string `xml:"info-name"`
	MaxSpeed              int    `xml:"max-speed"`
	MediaType             string `xml:"media-type"`
	MpiFirmwareRev        string `xml:"mpi-firmware-rev"`
	Node                  string `xml:"node"`
	NodeName              string `xml:"node-name"`
	PacketSize            int    `xml:"packet-size"`
	PciBusWidth           int    `xml:"pci-bus-width"`
	PciClockSpeed         int    `xml:"pci-clock-speed"`
	PhyFirmwareRev        string `xml:"phy-firmware-rev"`
	PhysicalDataLinkRate  int    `xml:"physical-data-link-rate"`
	PhysicalLinkState     string `xml:"physical-link-state"`
	PhysicalProtocol      string `xml:"physical-protocol"`
	PortAddress           int    `xml:"port-address"`
	PortName              string `xml:"port-name"`
	Speed                 string `xml:"speed"`
	SramParityEnabled     bool   `xml:"sram-parity-enabled"`
	State                 string `xml:"state"`
	SwitchPort            string `xml:"switch-port"`
	VlanId                int    `xml:"vlan-id"`
}

type FcpAdapterConfigGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			FcpAdapterAttributes []FcpAdapterConfigInfo `xml:"fcp-config-adapter-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type FcpAdapterConfigPageResponse struct {
	Response    *FcpAdapterConfigGetIterResponse
	Error       error
	RawResponse *http.Response
}

type FcpAdapterConfigPageHandler func(FcpAdapterConfigPageResponse) (shouldContinue bool)

func (f *Fcp) FcpAdapterGetIter(options *FcpAdapterConfigOptions) (*FcpAdapterConfigGetIterResponse, *http.Response, error) {
	f.Params.XMLName = xml.Name{Local: "fcp-adapter-get-iter"}
	f.Params.FcpAdapterConfigOptions = options
	r := FcpAdapterConfigGetIterResponse{}
	res, err := f.get(f, &r)
	return &r, res, err
}

func (f *Fcp) FcpAdapterGetAll(options *FcpAdapterConfigOptions, fn FcpAdapterConfigPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		fcpAdapterConfigGetIterResponse, res, err := f.FcpAdapterGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(FcpAdapterConfigPageResponse{Response: fcpAdapterConfigGetIterResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = fcpAdapterConfigGetIterResponse.Results.NextTag
			requestOptions = &FcpAdapterConfigOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
