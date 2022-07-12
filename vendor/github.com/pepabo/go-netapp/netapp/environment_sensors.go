package netapp

import (
	"encoding/xml"
	"net/http"
)

type EnvironmentSensors struct {
	Base
	Params struct {
		XMLName xml.Name
		*EnvironmentSensorsOptions
	}
}

func (e *EnvironmentSensors) List(options *EnvironmentSensorsOptions) (*EnvironmentSensorsResponse, *http.Response, error) {
	e.Params.XMLName = xml.Name{Local: "environment-sensors-get-iter"}
	e.Params.EnvironmentSensorsOptions = options
	r := EnvironmentSensorsResponse{}
	res, err := e.get(e, &r)
	return &r, res, err
}

func (e *EnvironmentSensors) ListPages(options *EnvironmentSensorsOptions, fn EnvironmentSensorsPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		response, res, err := e.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(EnvironmentSensorsPagesResponse{Response: response, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = response.Results.NextTag
			requestOptions = &EnvironmentSensorsOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type EnvironmentSensorsInfo struct {
	CriticalHighThreshold int    `xml:"critical-high-threshold,omitempty"`
	CriticalLowThreshold  int    `xml:"critical-low-threshold,omitempty"`
	DiscreteSensorState   string `xml:"discrete-sensor-state,omitempty"`
	DiscreteSensorValue   string `xml:"discrete-sensor-value,omitempty"`
	NodeName              string `xml:"node-name,omitempty"`
	SensorName            string `xml:"sensor-name,omitempty"`
	SensorType            string `xml:"sensor-type,omitempty"`
	ThresholdSensorState  string `xml:"threshold-sensor-state,omitempty"`
	ThresholdSensorValue  int    `xml:"threshold-sensor-value,omitempty"`
	ValueUnits            string `xml:"value-units,omitempty"`
	WarningHighThreshold  int    `xml:"warning-high-threshold,omitempty"`
	WarningLowThreshold   int    `xml:"warning-low-threshold,omitempty"`
}

type EnvironmentSensorsQuery struct {
	EnvironmentSensorsInfo *EnvironmentSensorsInfo `xml:"environment-sensors-info,omitempty"`
}

type EnvironmentSensorsOptions struct {
	DesiredAttributes *EnvironmentSensorsQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                      `xml:"max-records,omitempty"`
	Query             *EnvironmentSensorsQuery `xml:"query,omitempty"`
	Tag               string                   `xml:"tag,omitempty"`
}

type EnvironmentSensorsResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		EnvironmentSensorsInfo []EnvironmentSensorsInfo `xml:"attributes-list>environment-sensors-info"`
		NextTag                string                   `xml:"next-tag"`
		NumRecords             int                      `xml:"num-records"`
	} `xml:"results"`
}

type EnvironmentSensorsPagesResponse struct {
	Response    *EnvironmentSensorsResponse
	Error       error
	RawResponse *http.Response
}

type EnvironmentSensorsPageHandler func(EnvironmentSensorsPagesResponse) (shouldContinue bool)
