package netapp

import (
	"encoding/xml"
	"net/http"
)

type Snapshot struct {
	Base
	Params struct {
		XMLName xml.Name
		*SnapshotOptions
	}
}

type SnapshotQuery struct {
	SnapshotInfo *SnapshotInfo `xml:"snapshot-info,omitempty"`
}

type SnapshotOptions struct {
	DesiredAttributes *SnapshotQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int            `xml:"max-records,omitempty"`
	Query             *SnapshotQuery `xml:"query,omitempty"`
	Tag               string         `xml:"tag,omitempty"`
}

type SnapshotInfo struct {
	AccessTime                        int    `xml:"access-time"`
	Busy                              bool   `xml:"busy"`
	ContainsLunClones                 bool   `xml:"contains-lun-clones"`
	CumulativePercentageOfTotalBlocks int    `xml:"cumulative-percentage-of-total-blocks"`
	CumulativePercentageOfUsedBlocks  int    `xml:"cumulative-percentage-of-used-blocks"`
	CumulativeTotal                   int    `xml:"cumulative-total"`
	Dependency                        string `xml:"dependency"`
	Is7ModeSnapshot                   bool   `xml:"is-7-mode-snapshot"`
	Name                              string `xml:"name"`
	PercentageOfTotalBlocks           int    `xml:"percentage-of-total-blocks"`
	PercentageOfUsedBlocks            int    `xml:"percentage-of-used-blocks"`
	SnapmirrorLabel                   string `xml:"snapmirror-label"`
	SnapshotInstanceUuid              string `xml:"snapshot-instance-uuid"`
	SnapshotVersionUuid               string `xml:"snapshot-version-uuid"`
	State                             string `xml:"state"`
	Total                             int    `xml:"total"`
	Volume                            string `xml:"volume"`
	VolumeProvenanceUuid              string `xml:"volume-provenance-uuid"`
	Vserver                           string `xml:"vserver"`
}

type SnapshotListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			SnapshotAttributes []SnapshotInfo `xml:"snapshot-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type SnapshotListPagesResponse struct {
	Response    *SnapshotListResponse
	Error       error
	RawResponse *http.Response
}

type SnapshotPageHandler func(SnapshotListPagesResponse) (shouldContinue bool)

func (v *Snapshot) List(options *SnapshotOptions) (*SnapshotListResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: "snapshot-get-iter"}
	v.Params.SnapshotOptions = options
	r := SnapshotListResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}

func (v *Snapshot) ListPages(options *SnapshotOptions, fn SnapshotPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		snapshotResponse, res, err := v.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(SnapshotListPagesResponse{Response: snapshotResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = snapshotResponse.Results.NextTag
			requestOptions = &SnapshotOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
