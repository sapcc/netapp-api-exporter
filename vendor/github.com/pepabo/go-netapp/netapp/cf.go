package netapp

import (
	"encoding/xml"
	"net/http"
)

type Cf struct {
	Base
	Params struct {
		XMLName xml.Name
		*ClusterFailoverInfoOptions
	}
}

func (s *Cf) ClusterFailoverInfoList(options *ClusterFailoverInfoOptions) (*ClusterFailoverInfoResponse, *http.Response, error) {
	s.Params.XMLName = xml.Name{Local: "cf-get-iter"}
	s.Params.ClusterFailoverInfoOptions = options
	r := ClusterFailoverInfoResponse{}
	res, err := s.get(s, &r)
	return &r, res, err
}

func (s *Cf) ClusterFailoverInfoListPages(options *ClusterFailoverInfoOptions, fn StorageFailoverInfoPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		response, res, err := s.ClusterFailoverInfoList(requestOptions)
		handlerResponse := false

		handlerResponse = fn(ClusterFailoverInfoPagesResponse{Response: response, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = response.Results.NextTag
			requestOptions = &ClusterFailoverInfoOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type StorageFailoverInfoQuery struct {
	StorageFailoverInfo *StorageFailoverInfo `xml:"storage-failover-info,omitempty"`
}

type ClusterFailoverInfoOptions struct {
	DesiredAttributes *StorageFailoverInfoQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                       `xml:"max-records,omitempty"`
	Query             *StorageFailoverInfoQuery `xml:"query,omitempty"`
	Tag               string                    `xml:"tag,omitempty"`
}

type StorageFailoverInfo struct {
	InterconnectRelatedInfo *InterconnectRelatedInfo `xml:"sfo-interconnect-info>interconnect-related-info"`
	NodeRelatedInfo         *NodeRelatedInfo         `xml:"sfo-node-info>node-related-info"`
	TakeoverRelatedInfo     *TakeoverRelatedInfo     `xml:"sfo-takeover-info>takeover-related-info"`
	GivebackRelatedInfo     *GivebackRelatedInfo     `xml:"sfo-giveback-info>giveback-related-info"`
}

type InterconnectRelatedInfo struct {
	InterconnectLinks string `xml:"interconnect-links"`
	InterconnectType  string `xml:"interconnect-type"`
	IsInterconnectUp  bool   `xml:"is-interconnect-up"`
}

type NodeRelatedInfo struct {
	CurrentMode             string `xml:"current-mode"`
	LocalFirmwareProgress   int    `xml:"local-firmware-progress"`
	LocalFirmwareState      string `xml:"local-firmware-state"`
	Node                    string `xml:"node"`
	NodeState               string `xml:"node-state"`
	NvramId                 int    `xml:"nvram-id"`
	PartnerFirmwareProgress int    `xml:"partner-firmware-progress"`
	PartnerFirmwareState    string `xml:"partner-firmware-state"`
	PartnerName             string `xml:"partner-name"`
	PartnerNvramId          int    `xml:"partner-nvram-id"`
	StateDescription        string `xml:"state-description"`
}

type TakeoverRelatedInfo struct {
	TakeoverByPartnerPossible bool   `xml:"takeover-by-partner-possible"`
	TakeoverEnabled           bool   `xml:"takeover-enabled"`
	TakeoverFailureReason     string `xml:"takeover-failure-reason"`
	TakeoverModule            string `xml:"takeover-module"`
	TakeoverOfPartnerPossible bool   `xml:"takeover-of-partner-possible"`
	TakeoverReason            string `xml:"takeover-reason"`
	TakeoverState             string `xml:"takeover-state"`
	TimeSinceTakeover         int    `xml:"time-since-takeover"`
	TimeUntilTakeover         int    `xml:"time-until-takeover"`
}

type GivebackRelatedInfo struct {
	GivebackState string `xml:"giveback-state"`
}

type ClusterFailoverInfoResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList []StorageFailoverInfo `xml:"attributes-list>storage-failover-info"`
		NextTag        string                `xml:"next-tag"`
		NumRecords     int                   `xml:"num-records"`
	} `xml:"results"`
}

type ClusterFailoverInfoPagesResponse struct {
	Response    *ClusterFailoverInfoResponse
	Error       error
	RawResponse *http.Response
}

type StorageFailoverInfoPageHandler func(ClusterFailoverInfoPagesResponse) (shouldContinue bool)
