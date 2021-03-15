package netapp

import (
	"encoding/xml"
	"net/http"
)

type Diagnosis struct {
	Base
	Params struct {
		XMLName xml.Name
		*DiagnosisOptions
	}
}

type DiagnosisQuery struct {
	DiagnosisInfo *DiagnosisAlertInfo `xml:"diagnosis-alert-info,omitempty"`
}

type DiagnosisOptions struct {
	DesiredAttributes *DiagnosisAlertInfo `xml:"desired-attributes,omitempty"`
	MaxRecords        int                 `xml:"max-records,omitempty"`
	Query             *DiagnosisQuery     `xml:"query,omitempty"`
	Tag               string              `xml:"tag,omitempty"`
}

type DiagnosisAlertInfo struct {
	Acknowledge              bool   `xml:"acknowledge"`
	Acknowledger             string `xml:"acknowledger"`
	Additionalinfo           string `xml:"additional-info"`
	AlertId                  string `xml:"alert-id"`
	AlertingResource         string `xml:"alerting-resource"`
	AlertingResourceName     string `xml:"alerting-resource-name"`
	CorrectiveActions        string `xml:"corrective-actions"`
	IndicationTime           int    `xml:"indication-time"`
	Monitor                  string `xml:"monitor"`
	Node                     string `xml:"node"`
	PerceivedSeverity        string `xml:"perceived-severity"`
	Policy                   string `xml:"policy"`
	PossibleEffect           string `xml:"possible-effect"`
	ProbableCause            string `xml:"probable-cause"`
	ProbableCauseDescription string `xml:"probable-cause-description"`
	Subsystem                string `xml:"subsystem"`
	Suppress                 bool   `xml:"suppress"`
	Suppressor               string `xml:"suppressor"`
}

type DiagnosisListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			DiagnosisAttributes []DiagnosisAlertInfo `xml:"diagnosis-alert-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type DiagnosisAlertPagesResponse struct {
	Response    *DiagnosisListResponse
	Error       error
	RawResponse *http.Response
}

type DiagnosisPageHandler func(DiagnosisAlertPagesResponse) (shouldContinue bool)

func (v *Diagnosis) DiagnosisAlertGetIter(options *DiagnosisOptions) (*DiagnosisListResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: "diagnosis-alert-get-iter"}
	v.Params.DiagnosisOptions = options
	r := DiagnosisListResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}

func (v *Diagnosis) DiagnosisAlertGetAll(options *DiagnosisOptions, fn DiagnosisPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		DiagnosisResponse, res, err := v.DiagnosisAlertGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(DiagnosisAlertPagesResponse{Response: DiagnosisResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = DiagnosisResponse.Results.NextTag
			requestOptions = &DiagnosisOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
