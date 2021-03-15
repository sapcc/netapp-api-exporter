package netapp

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Qtree struct {
	Base
	Params struct {
		XMLName xml.Name
		*QtreeOptions
	}
}

type QtreeQuery struct {
	QtreeInfo *QtreeInfo `xml:"qtree-info,omitempty"`
}

type QtreeOptions struct {
	DesiredAttributes *QtreeQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int         `xml:"max-records,omitempty"`
	Query             *QtreeQuery `xml:"query,omitempty"`
	Tag               string      `xml:"tag,omitempty"`
	*QtreeInfo
}

type QtreeInfo struct {
	ExportPolicy            string `xml:"export-policy,omitempty"`
	ID                      string `xml:"id,omitempty"`
	IsExportPolicyInherited string `xml:"is-export-policy-inherited,omitempty"`
	Mode                    string `xml:"mode,omitempty"`
	Oplocks                 string `xml:"oplocks,omitempty"`
	Qtree                   string `xml:"qtree,omitempty"`
	Force                   bool   `xml:"force,omitempty"`
	SecurityStyle           string `xml:"security-style,omitempty"`
	Status                  string `xml:"status,omitempty"`
	Volume                  string `xml:"volume,omitempty"`
	Vserver                 string `xml:"vserver,omitempty"`
}

type QtreeListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			QtreeInfo []QtreeInfo `xml:"qtree-info"`
		} `xml:"attributes-list"`
	} `xml:"results"`
	ResultJobid  string `xml:"result-jobid"`
	ResultStatus string `xml:"result-status"`
}

func (q *Qtree) List(options *QtreeOptions) (*QtreeListResponse, *http.Response, error) {
	q.Params.XMLName = xml.Name{Local: "qtree-list-iter"}
	q.Params.QtreeOptions = options

	r := QtreeListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Qtree) Create(vserverName, volume, qtree string, info *QtreeInfo) (*QtreeListResponse, *http.Response, error) {
	q.Name = vserverName
	q.Params.XMLName = xml.Name{Local: "qtree-create"}
	if info == nil {
		info = &QtreeInfo{}
	}
	info.Volume = volume
	info.Qtree = qtree

	q.Params.QtreeOptions = &QtreeOptions{
		QtreeInfo: info,
	}

	r := QtreeListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Qtree) Delete(vserverName, volName, qtreeName string, force bool) (*QtreeListResponse, *http.Response, error) {
	q.Name = vserverName
	q.Params.XMLName = xml.Name{Local: "qtree-delete"}
	q.Params.QtreeOptions = &QtreeOptions{
		QtreeInfo: &QtreeInfo{
			Qtree: fmt.Sprintf("/vol/%s/%s", volName, qtreeName),
			Force: force,
		},
	}

	r := QtreeListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}

func (q *Qtree) DeleteAsync(vserverName, volName, qtreeName string) (*QtreeListResponse, *http.Response, error) {
	q.Name = vserverName
	q.Params.XMLName = xml.Name{Local: "qtree-delete-async"}
	q.Params.QtreeOptions = &QtreeOptions{
		QtreeInfo: &QtreeInfo{
			Qtree: fmt.Sprintf("/vol/%s/%s", volName, qtreeName),
		},
	}

	r := QtreeListResponse{}
	res, err := q.get(q, &r)
	return &r, res, err
}
