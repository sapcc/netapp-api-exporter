package netapp

import (
	"encoding/xml"
	"net/http"
)

const (
	QuotaStatusCorrupt      = "corrupt"
	QuotaStatusInitializing = "initializing"
	QuotaStatusMixed        = "mixed"
	QuotaStatusOff          = "off"
	QuotaStatusOn           = "on"
	QuotaStatusResizing     = "resizing"
	QuotaStatusReverting    = "reverting"
	QuotaStatusUnknown      = "unknown"
	QuotaStatusUpgrading    = "upgrading"
)

type Quota struct {
	Base
	Params struct {
		XMLName xml.Name
		*QuotaOptions
	}
}

type QuotaQuery struct {
	QuotaEntry *QuotaEntry `xml:"quota-entry,omitempty"`
}

type QuotaOptions struct {
	DesiredAttributes *QuotaQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int         `xml:"max-records,omitempty"`
	Query             *QuotaQuery `xml:"query,omitempty"`
	Tag               string      `xml:"tag,omitempty"`
	*QuotaEntry
}

type QuotaEntry struct {
	DiskLimit          string  `xml:"disk-limit,omitempty"`
	FileLimit          string  `xml:"file-limit,omitempty"`
	PerformUserMapping string  `xml:"perform-user-mapping,omitempty"`
	Policy             string  `xml:"policy,omitempty"`
	Qtree              *string `xml:"qtree,omitempty"`
	QuotaTarget        string  `xml:"quota-target,omitempty"`
	QuotaType          string  `xml:"quota-type,omitempty"`
	SoftDiskLimit      string  `xml:"soft-disk-limit,omitempty"`
	SoftFileLimit      string  `xml:"soft-file-limit,omitempty"`
	Threshold          string  `xml:"threshold,omitempty"`
	Volume             string  `xml:"volume,omitempty"`
	Vserver            string  `xml:"vserver,omitempty"`
}

type QuotaResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		QuotaEntry
	} `xml:"results"`
}

type QuotaListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			QuotaEntry []QuotaEntry `xml:"quota-entry"`
		} `xml:"attributes-list"`
	} `xml:"results"`
}

type QuotaStatusResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		QuotaStatus    string `xml:"status"`
		QuotaSubStatus string `xml:"substatus"`
		ResultJobid    string `xml:"result-jobid"`
		ResultStatus   string `xml:"result-status"`
	} `xml:"results"`
}

func (q *Quota) Get(name string, options *QuotaOptions) (*QuotaResponse, *http.Response, error) {
	q.Name = name
	q.Params.XMLName = xml.Name{Local: "quota-get-entry"}
	q.Params.QuotaOptions = options
	r := QuotaResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) List(options *QuotaOptions) (*QuotaListResponse, *http.Response, error) {
	q.Params.XMLName = xml.Name{Local: "quota-list-entries-iter"}
	q.Params.QuotaOptions = options

	r := QuotaListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) Create(serverName, target, quotaType, qtree string, entry *QuotaEntry) (*QuotaListResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-add-entry"}

	if entry == nil {
		entry = &QuotaEntry{}
	}

	entry.QuotaTarget = target
	entry.QuotaType = quotaType
	entry.Qtree = &qtree

	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: entry,
	}

	r := QuotaListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) Update(serverName string, entry *QuotaEntry) (*QuotaListResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-modify-entry"}
	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: entry,
	}

	r := QuotaListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) Delete(serverName, target, quotaType, volume, qtree string) (*QuotaListResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-delete-entry"}
	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: &QuotaEntry{
			QuotaType:   quotaType,
			QuotaTarget: target,
			Volume:      volume,
			Qtree:       &qtree,
		},
	}

	r := QuotaListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) Off(serverName, volumeName string) (*QuotaStatusResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-off"}
	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: &QuotaEntry{
			Volume: volumeName,
		},
	}

	r := QuotaStatusResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) On(serverName, volumeName string) (*QuotaStatusResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-on"}
	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: &QuotaEntry{
			Volume: volumeName,
		},
	}

	r := QuotaStatusResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Quota) Status(serverName, volumeName string) (*QuotaStatusResponse, *http.Response, error) {
	q.Name = serverName
	q.Params.XMLName = xml.Name{Local: "quota-status"}
	q.Params.QuotaOptions = &QuotaOptions{
		QuotaEntry: &QuotaEntry{
			Volume: volumeName,
		},
	}

	r := QuotaStatusResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

type QuotaReport struct {
	Base
	Params struct {
		XMLName xml.Name
		*QuotaReportOptions
	}
}

func (qr *QuotaReport) Report(options *QuotaReportOptions) (*QuotaReportResponse, *http.Response, error) {
	qr.Params.XMLName = xml.Name{Local: "quota-report-iter"}
	qr.Params.QuotaReportOptions = options

	r := QuotaReportResponse{}
	res, err := qr.get(qr, &r)
	return &r, res, err
}

type QuotaReportPageHandler func(QuotaReportPagesResponse) (shouldContinue bool)

func (a *QuotaReport) ReportPages(options *QuotaReportOptions, fn QuotaReportPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		quotaReportResponse, res, err := a.Report(requestOptions)
		handlerResponse := false

		handlerResponse = fn(QuotaReportPagesResponse{Response: quotaReportResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = quotaReportResponse.Results.NextTag
			requestOptions = &QuotaReportOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type QuotaReportEntryQuery struct {
	QuotaReportEntry *QuotaReportEntry `xml:"quota,omitempty"`
}

type QuotaReportOptions struct {
	DesiredAttributes *QuotaReportEntryQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                    `xml:"max-records,omitempty"`
	Path              string                 `xml:"path,omitempty"`
	Query             *QuotaReportEntryQuery `xml:"query,omitempty"`
	Tag               string                 `xml:"tag,omitempty"`
}

type QuotaReportEntry struct {
	DiskLimit     string `xml:"disk-limit,omitempty"`
	DiskUsed      string `xml:"disk-used,omitempty"`
	FileLimit     string `xml:"file-limit,omitempty"`
	FilesUsed     string `xml:"files-used,omitempty"`
	QuotaTarget   string `xml:"quota-target,omitempty"`
	QuotaType     string `xml:"quota-type,omitempty"`
	SoftDiskLimit string `xml:"soft-disk-limit,omitempty"`
	SoftFileLimit string `xml:"soft-file-limit,omitempty"`
	Threshold     string `xml:"threshold,omitempty"`
	Tree          string `xml:"tree,omitempty"`
	Volume        string `xml:"volume,omitempty"`
	Vserver       string `xml:"vserver,omitempty"`
}

type QuotaReportResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			QuotaReportEntry []QuotaReportEntry `xml:"quota"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type QuotaReportPagesResponse struct {
	Response    *QuotaReportResponse
	Error       error
	RawResponse *http.Response
}

type QuotaStatus struct {
	Base
	Params struct {
		XMLName xml.Name
		*QuotaStatusIterOptions
	}
}

func (qr *QuotaStatus) StatusIter(options *QuotaStatusIterOptions) (*QuotaStatusIterResponse, *http.Response, error) {
	qr.Params.XMLName = xml.Name{Local: "quota-status-iter"}
	qr.Params.QuotaStatusIterOptions = options

	r := QuotaStatusIterResponse{}
	res, err := qr.get(qr, &r)
	return &r, res, err
}

func (a *QuotaStatus) StatusPages(options *QuotaStatusIterOptions, fn QuotaStatusPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		quotaStatusResponse, res, err := a.StatusIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(QuotaStatusPagesResponse{Response: quotaStatusResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = quotaStatusResponse.Results.NextTag
			requestOptions = &QuotaStatusIterOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type QuotaStatusEntryQuery struct {
	QuotaStatusEntry *QuotaStatusEntry `xml:"quota-status-attributes,omitempty"`
}

type QuotaStatusIterOptions struct {
	DesiredAttributes *QuotaStatusEntryQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                    `xml:"max-records,omitempty"`
	Query             *QuotaStatusEntryQuery `xml:"query,omitempty"`
	Tag               string                 `xml:"tag,omitempty"`
}

type QuotaStatusEntry struct {
	PercentComplete string `xml:"percent-complete"`
	QuotaErrorMsgs  string `xml:"quota-error-msgs"`
	Reason          string `xml:"reason"`
	QuotaStatus     string `xml:"status"`
	QuotaSubStatus  string `xml:"substatus"`
	Volume          string `xml:"volume"`
	Vserver         string `xml:"vserver"`
}

type QuotaStatusIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			QuotaStatusAttributes []QuotaStatusEntry `xml:"quota-status-attributes"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type QuotaStatusPagesResponse struct {
	Response    *QuotaStatusIterResponse
	Error       error
	RawResponse *http.Response
}

type QuotaStatusPageHandler func(QuotaStatusPagesResponse) (shouldContinue bool)
