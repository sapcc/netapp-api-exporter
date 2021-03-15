package netapp

import (
	"encoding/xml"
	"net/http"
)

type System struct {
	Base
	Params struct {
		XMLName xml.Name
		*NodeDetailOptions
	}
}

func (s *System) List(options *NodeDetailOptions) (*NodeDetailsResponse, *http.Response, error) {
	s.Params.XMLName = xml.Name{Local: "system-node-get-iter"}
	s.Params.NodeDetailOptions = options
	r := NodeDetailsResponse{}
	res, err := s.get(s, &r)
	return &r, res, err
}

func (s *System) ListPages(options *NodeDetailOptions, fn NodeDetailsPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		response, res, err := s.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(NodeDetailsPagesResponse{Response: response, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = response.Results.NextTag
			requestOptions = &NodeDetailOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type NodeDetails struct {
	EnvFailedFanCount           int    `xml:"env-failed-fan-count"`
	EnvFailedFanMessage         string `xml:"env-failed-fan-message"`
	EnvFailedPowerSupplyCount   int    `xml:"env-failed-power-supply-count"`
	EnvFailedPowerSupplyMessage string `xml:"env-failed-power-supply-message"`
	EnvOverTemperature          bool   `xml:"env-over-temperature"`
	Name                        string `xml:"node"`
	NodeAssetTag                string `xml:"node-asset-tag"`
	NodeLocation                string `xml:"node-location"`
	NodeModel                   string `xml:"node-model"`
	NodeNvramId                 string `xml:"node-nvram-id"`
	NodeOwner                   string `xml:"node-owner"`
	NodeSerialNumber            string `xml:"node-serial-number"`
	NodeStorageConfiguration    string `xml:"node-storage-configuration"`
	NodeSystemId                string `xml:"node-system-id"`
	NodeUptime                  string `xml:"node-uptime"`
	NodeUuid                    string `xml:"node-uuid"`
	NodeVendor                  string `xml:"node-vendor"`
	NvramBatteryStatus          string `xml:"nvram-battery-status"`
	ProductVersion              string `xml:"product-version"`
}

type NodeDetailsQuery struct {
	NodeDetails *NodeDetails `xml:"node-details-info,omitempty"`
}

type NodeDetailOptions struct {
	DesiredAttributes *NodeDetailsQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int               `xml:"max-records,omitempty"`
	Query             *NodeDetailsQuery `xml:"query,omitempty"`
	Tag               string            `xml:"tag,omitempty"`
}

type NodeDetailsResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		NodeDetails []NodeDetails `xml:"attributes-list>node-details-info"`
		NextTag     string        `xml:"next-tag"`
		NumRecords  int           `xml:"num-records"`
	} `xml:"results"`
}

type NodeDetailsPagesResponse struct {
	Response    *NodeDetailsResponse
	Error       error
	RawResponse *http.Response
}

type NodeDetailsPageHandler func(NodeDetailsPagesResponse) (shouldContinue bool)
