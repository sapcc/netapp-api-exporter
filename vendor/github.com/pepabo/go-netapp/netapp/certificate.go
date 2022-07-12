package netapp

import (
	"encoding/xml"
	"net/http"
)

type Certificate struct {
	Base
	Params struct {
		XMLName xml.Name
		*CertificateOptions
	}
}

type CertificateQuery struct {
	CertificateInfo *CertificateInfo `xml:"certificate-info,omitempty"`
}

type CertificateOptions struct {
	DesiredAttributes *CertificateQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int               `xml:"max-records,omitempty"`
	Query             *CertificateQuery `xml:"query,omitempty"`
	Tag               string            `xml:"tag,omitempty"`
}

type CertificateInfo struct {
	Vserver        string `xml:"vserver,omitempty"`
	Type           string `xml:"type,omitempty"`
	SerialNumber   string `xml:"serial-number,omitempty"`
	CommonName     string `xml:"common-name,omitempty"`
	ExpirationDate int    `xml:"expiration-date,omitempty"`
}

type CertificateListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList []CertificateInfo `xml:"attributes-list>certificate-info"`
		NextTag        string            `xml:"next-tag"`
		NumRecords     int               `xml:"num-records"`
	} `xml:"results"`
}

type CertificatePagesResponse struct {
	Response    *CertificateListResponse
	Error       error
	RawResponse *http.Response
}

type CertificatePageHandler func(CertificatePagesResponse) (shouldContinue bool)

func (v *Certificate) CertificateGetIter(options *CertificateOptions) (*CertificateListResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: "security-certificate-get-iter"}
	v.Params.CertificateOptions = options
	r := CertificateListResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}

func (v *Certificate) CertificateGetAll(options *CertificateOptions, fn CertificatePageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		CertificateResponse, res, err := v.CertificateGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(CertificatePagesResponse{Response: CertificateResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = CertificateResponse.Results.NextTag
			requestOptions = &CertificateOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
