package netapp

import (
	"encoding/xml"
	"net/http"
)

type VolumeModifyOptions struct {
	*VolumeOptions
	ContinueOnFailure bool `xml:"continue-on-failure,omitempty"`
	MaxFailureCount   int  `xml:"max-failure-count,omitempty"`
	ReturnFailureList bool `xml:"return-failure-list,omitempty"`
	ReturnSuccessList bool `xml:"return-success-list,omitempty"`
}

type VolumeModifyResponce struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		FailureList     *[]VolumeModifyInfo `xml:"failure-list>volume-modify-iter-info"`
		SuccessList     *[]VolumeModifyInfo `xml:"success-list>volume-modify-iter-info"`
		NextTag         string              `xml:"next-tag"`
		NumberSucceeded int                 `xml:"num-succeeded"`
		NumberFailed    int                 `xml:"num-failed"`
	} `xml:"results"`
}

type VolumeModifyInfo struct {
	ErrorCode    int         `xml:"error-code,omitempty"`
	ErrorMessage string      `xml:"error-message,omitempty"`
	VolumeKey    *VolumeInfo `xml:"volume-key,omitempty"`
}

// Modify changes some volume properties, note: it will silently ignore things it cannot change
func (v Volume) Modify(options *VolumeOptions) (*VolumeModifyResponce, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: "volume-modify-iter"}
	v.Params.VolumeOptions = options
	r := VolumeModifyResponce{}
	res, err := v.get(v, &r)
	return &r, res, err
}
