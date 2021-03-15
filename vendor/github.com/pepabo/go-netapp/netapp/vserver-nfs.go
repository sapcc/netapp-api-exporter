package netapp

import (
	"encoding/xml"
	"net/http"
)

type vServerNfsRequest struct {
	Base
	Params struct {
		XMLName                 xml.Name
		VServerNfsCreateOptions `xml:",innerxml"`
	}
}

type VServerNfsCreateOptions struct {
	NfsAccessEnabled bool `xml:"is-nfs-access-enabled"`
	NfsV3Enabled     bool `xml:"is-nfsv3-enabled"`
	NfsV4Enabled     bool `xml:"is-nfsv40-enabled"`
	VStorageEnabled  bool `xml:"is-vstorage-enabled"`
}

// CreateNfsService configures and enables nfs service on a vserver
func (v VServer) CreateNfsService(vServerName string, options *VServerNfsCreateOptions) (*SingleResultResponse, *http.Response, error) {
	req := v.newVServerNfsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "nfs-service-create"}
	req.Params.VServerNfsCreateOptions = *options

	r := &SingleResultResponse{}
	res, err := v.get(req, r)
	return r, res, err
}

func (v VServer) newVServerNfsRequest() *vServerNfsRequest {
	return &vServerNfsRequest{
		Base: v.Base,
	}
}
