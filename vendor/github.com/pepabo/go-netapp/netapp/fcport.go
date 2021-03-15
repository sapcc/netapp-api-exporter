package netapp

import (
	"encoding/xml"
	"net/http"
)

type Fcport struct {
	Base
	Params struct {
		XMLName xml.Name
		*FcportGetLinkStateOptions
	}
}

type FcportGetLinkStateOptions struct {
	AdapterName string `xml:"adapter-name,omitempty"`
	NodeName    string `xml:"node-name,omitempty"`
}

type FcportLinkStateInfo struct {
	AdapterName string `xml:"adapter-name"`
	LinkState   string `xml:"link-state"`
	NodeName    string `xml:"node-name"`
}

type FcportGetLinkStateResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AdapterLinkState []FcportLinkStateInfo `xml:"adapter-link-state"`
	} `xml:"results"`
}

func (f *Fcport) GetLinkState(options *FcportGetLinkStateOptions) (*FcportGetLinkStateResponse, *http.Response, error) {
	f.Params.XMLName = xml.Name{Local: "fcport-get-link-state"}
	f.Params.FcportGetLinkStateOptions = options
	r := FcportGetLinkStateResponse{}
	res, err := f.get(f, &r)
	return &r, res, err
}
