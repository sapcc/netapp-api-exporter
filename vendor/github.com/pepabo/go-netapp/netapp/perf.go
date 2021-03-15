package netapp

import (
	"encoding/xml"
	"net/http"
)

type Perf struct {
	Base
}

type PerfCounterData struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type PerfCounters struct {
	CounterData []PerfCounterData `xml:"counter-data"`
}

type InstanceData struct {
	Name     string       `xml:"name"`
	Counters PerfCounters `xml:"counters"`
}

type PerfObjectInstanceData struct {
	Instances []InstanceData `xml:"instance-data"`
}

type PerfObjectGetInstancesResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		PerfObjectInstanceData PerfObjectInstanceData `xml:"instances"`
	} `xml:"results"`
}

type PerfObjectGetInstanceParams struct {
	ObjectName    string `xml:"objectname"`
	InstanceUuids struct {
		Uuids []string `xml:"instance-uuid"`
	} `xml:"instance-uuids,omitempty"`
	Instances struct {
		Instances []string `xml:"instance"`
	} `xml:"instances,omitempty"`
	Counters struct {
		Counter []string `xml:"counter"`
	} `xml:"counters"`
}

type perfObjectGetInstanceRequest struct {
	Base
	Params struct {
		XMLName xml.Name
		*PerfObjectGetInstanceParams
	}
}

func newPerfObjectGetInstanceRequest(params *PerfObjectGetInstanceParams, base Base) *perfObjectGetInstanceRequest {
	request := perfObjectGetInstanceRequest{
		Base: base,
	}
	request.Params.XMLName = xml.Name{Local: "perf-object-get-instances"}
	request.Params.PerfObjectGetInstanceParams = params
	return &request
}

type InstanceInfoQuery struct {
	InstanceInfo *InstanceInfo `xml:"instance-info,omitempty"`
}

type PerfObjectInstanceListInfoIterParams struct {
	DesiredAttributes *InstanceInfo      `xml:"desired-attributes,omitempty"`
	FilterData        string             `xml:"filter-data,omitempty"`
	MaxRecords        int                `xml:"max-records,omitempty"`
	ObjectName        string             `xml:"objectname"`
	Query             *InstanceInfoQuery `xml:"query,omitempty"`
	Tag               string             `xml:"tag,omitempty"`
}

type InstanceInfo struct {
	Name string `xml:"name"`
	Uuid string `xml:"uuid"`
}

type PerfObjectInstanceListInfoIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			InstanceInfo []InstanceInfo `xml:"instance-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type PerfObjectInstanceListInfoPageResponse struct {
	Response    *PerfObjectInstanceListInfoIterResponse
	Error       error
	RawResponse *http.Response
}

type perfObjectInstanceListInfoIterRequest struct {
	Base
	Params struct {
		XMLName xml.Name
		*PerfObjectInstanceListInfoIterParams
	}
}

func newPerfObjectInstanceListInfoIterRequest(params *PerfObjectInstanceListInfoIterParams, base Base) *perfObjectInstanceListInfoIterRequest {
	request := perfObjectInstanceListInfoIterRequest{
		Base: base,
	}
	request.Params.XMLName = xml.Name{Local: "perf-object-instance-list-info-iter"}
	request.Params.PerfObjectInstanceListInfoIterParams = params
	return &request
}

func (p *Perf) PerfObjectGetInstances(params *PerfObjectGetInstanceParams) (*PerfObjectGetInstancesResponse, *http.Response, error) {
	request := newPerfObjectGetInstanceRequest(params, p.Base)
	response := PerfObjectGetInstancesResponse{}
	rawResponse, err := p.get(request, &response)
	return &response, rawResponse, err
}

func (p *Perf) PerfObjectInstanceListInfoIter(params *PerfObjectInstanceListInfoIterParams) (*PerfObjectInstanceListInfoIterResponse, *http.Response, error) {
	request := newPerfObjectInstanceListInfoIterRequest(params, p.Base)
	response := PerfObjectInstanceListInfoIterResponse{}
	rawResponse, err := p.get(request, &response)
	return &response, rawResponse, err
}

type PerfObjectInstanceListInfoHandler func(PerfObjectInstanceListInfoPageResponse) (shouldContinue bool)

func (p *Perf) PerfObjectInstanceGetAllInfo(options *PerfObjectInstanceListInfoIterParams, fn PerfObjectInstanceListInfoHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		response, rawResponse, err := p.PerfObjectInstanceListInfoIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(PerfObjectInstanceListInfoPageResponse{Response: response, Error: err, RawResponse: rawResponse})

		nextTag := ""
		if err == nil {
			nextTag = response.Results.NextTag
			requestOptions = &PerfObjectInstanceListInfoIterParams{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}
