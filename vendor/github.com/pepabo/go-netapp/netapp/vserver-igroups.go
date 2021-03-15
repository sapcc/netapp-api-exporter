package netapp

import (
	"encoding/xml"
	"net/http"
)

type vServerIgroupsRequest struct {
	Base
	Params struct {
		XMLName           xml.Name
		VServerIgroupInfo `xml:",innerxml"`
	}
}

type vServerIgroupInfo struct {
	Base
	Params struct {
		XMLName           xml.Name
		VServerIgroupInfo `xml:",innerxml"`
	}
}

// VServerIgroupInfo sets all different options for the igroup
type VServerIgroupInfo struct {
	PortsetName        string `xml:"initiator-group-portset-name,omitempty"`
	InitiatorGroupName string `xml:"initiator-group-name,omitempty"`
	InitiatorGroupUUID string `xml:"initiator-group-uuid,omitempty"`
	InitiatorName      string `xml:"initiator,omitempty"`
	Force              bool   `xml:"force,omitempty"`
}

// VServerIgroupsResponse creates correct response obj
type VServerIgroupsResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
	} `xml:"results"`
}

// AddInitiator add an initiator to an igroup
func (v VServer) AddInitiator(vServerName string, iGroupName string, initiator string) (*VServerIgroupsResponse, *http.Response, error) {
	req := v.newVServerIgroupsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "igroup-add"}
	req.Params.VServerIgroupInfo = VServerIgroupInfo{
		InitiatorGroupName: iGroupName,
		InitiatorName:      initiator,
	}

	r := &VServerIgroupsResponse{}
	res, err := v.get(req, r)
	return r, res, err
}

// RemoveInitiator add an initiator to an igroup
func (v VServer) RemoveInitiator(vServerName string, iGroupName string, initiator string, force bool) (*VServerIgroupsResponse, *http.Response, error) {
	req := v.newVServerIgroupsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "igroup-remove"}
	req.Params.VServerIgroupInfo = VServerIgroupInfo{
		InitiatorGroupName: iGroupName,
		InitiatorName:      initiator,
		Force:              force,
	}

	r := &VServerIgroupsResponse{}
	res, err := v.get(req, r)
	return r, res, err

}

func (v VServer) newVServerIgroupsRequest() *vServerIgroupsRequest {
	return &vServerIgroupsRequest{
		Base: v.Base,
	}
}
