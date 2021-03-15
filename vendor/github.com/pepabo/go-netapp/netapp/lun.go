package netapp

import (
	"encoding/xml"
	"net/http"
)

type Lun struct {
	Base
	Params struct {
		XMLName xml.Name
		*LunOptions
	}
}
type LunQuery struct {
	LunInfo *LunInfo `xml:"lun-info,omitempty"`
}

type LunOptions struct {
	DesiredAttributes *LunQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int       `xml:"max-records,omitempty"`
	Query             *LunQuery `xml:"query,omitempty"`
	Tag               string    `xml:"tag,omitempty"`
}

type LunInfo struct {
	Alignment                 string `xml:"alignment"`
	BackingSnapshot           string `xml:"backing-snapshot"`
	BlockSize                 int    `xml:"block-size"`
	Class                     string `xml:"class"`
	CloneBackingSnapshot      string `xml:"clone-backing-snapshot"`
	Comment                   string `xml:"comment"`
	CreationTimestamp         int    `xml:"creation-timestamp"`
	DeviceBinaryId            string `xml:"device-binary-id"`
	DeviceId                  int    `xml:"device-id"`
	DeviceTextId              string `xml:"device-text-id"`
	IsClone                   bool   `xml:"is-clone"`
	IsCloneAutodeleteEnabled  bool   `xml:"is-clone-autodelete-enabled"`
	IsInconsistentImport      bool   `xml:"is-inconsistent-import"`
	IsRestoreInaccessible     bool   `xml:"is-restore-inaccessible"`
	IsSpaceAllocEnabled       bool   `xml:"is-space-alloc-enabled"`
	IsSpaceReservationEnabled bool   `xml:"is-space-reservation-enabled"`
	Mapped                    bool   `xml:"mapped"`
	MultiprotocolType         string `xml:"multiprotocol-type"`
	Node                      string `xml:"node"`
	Online                    bool   `xml:"online"`
	Path                      string `xml:"path"`
	PrefixSize                int    `xml:"prefix-size"`
	QosPolicyGroup            string `xml:"qos-policy-group"`
	Qtree                     string `xml:"qtree"`
	ReadOnly                  bool   `xml:"read-only"`
	Serial7Mode               string `xml:"serial-7-mode"`
	SerialNumber              string `xml:"serial-number"`
	ShareState                string `xml:"share-state"`
	Size                      int    `xml:"size"`
	SizeUsed                  int    `xml:"size-used"`
	Staging                   bool   `xml:"staging"`
	State                     string `xml:"state"`
	SuffixSize                int    `xml:"suffix-size"`
	Uuid                      string `xml:"uuid"`
	Volume                    string `xml:"volume"`
	Vserver                   string `xml:"vserver"`
}

type LunListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			LunAttributes []LunInfo `xml:"lun-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type LunListPagesResponse struct {
	Response    *LunListResponse
	Error       error
	RawResponse *http.Response
}

type LunPageHandler func(LunListPagesResponse) (shouldContinue bool)

func (v *Lun) List(options *LunOptions) (*LunListResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: "lun-get-iter"}
	v.Params.LunOptions = options
	r := LunListResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}

func (v *Lun) ListPages(options *LunOptions, fn LunPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		LunResponse, res, err := v.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(LunListPagesResponse{Response: LunResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = LunResponse.Results.NextTag
			requestOptions = &LunOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
