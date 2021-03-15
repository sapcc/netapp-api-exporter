package netapp

import (
	"encoding/xml"
	"net/http"
)

// QosPolicy is the main struct we're building on
type QosPolicy struct {
	Base
	Params struct {
		XMLName xml.Name
		Query   *QosPolicyInfo `xml:"desired-attributes,omitempty"`
		QosPolicyInfo
		QosPolicyRenameInfo
	}
}

// QosPolicyInfo is all qos policy data netapp stores
type QosPolicyInfo struct {
	MaxThroughput    string `xml:"max-throughput,omitempty"`
	NumWorkloads     int    `xml:"num-workloads,omitempty"`
	PgID             int    `xml:"pgid,omitempty"`
	PolicyGroup      string `xml:"policy-group,omitempty"`
	PolicyGroupClass string `xml:"policy-group-class,omitempty"`
	UUID             string `xml:"uuid,omitempty"`
	VServer          string `xml:"vserver,omitempty"`
	// ReturnRecord is only used in create
	ReturnRecord bool `xml:"return-record,omitempty"`
	// Force is only used in delete
	Force bool `xml:"force,omitempty"`
}

// QosPolicyRenameInfo is a struct for renaming a qos policy
type QosPolicyRenameInfo struct {
	CurrentPolicyGroup string `xml:"policy-group-name,omitempty"`
	NewPolicyGroup     string `xml:"new-name,omitempty"`
}

// QosPolicyResponse is what comes back from the api
type QosPolicyResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		QosPolicyInfo QosPolicyInfo `xml:"attributes>qos-policy-group-info"`
	} `xml:"results"`
}

// Create makes new qos policy
func (qp QosPolicy) Create(query *QosPolicyInfo) (*QosPolicyResponse, *http.Response, error) {
	qp.Params.XMLName = xml.Name{Local: "qos-policy-group-create"}
	qp.Params.QosPolicyInfo = *query

	return qp.doAPICall()
}

// Get grabs a qos policy, note: it will do so cluster wide
func (qp QosPolicy) Get(name string, query *QosPolicyInfo) (*QosPolicyResponse, *http.Response, error) {
	qp.Params.XMLName = xml.Name{Local: "qos-policy-group-get"}
	qp.Params.Query = query
	qp.Params.QosPolicyInfo.PolicyGroup = name

	return qp.doAPICall()
}

// Rename changes policy name, any volumes attached to the policy get the new name automatically
func (qp QosPolicy) Rename(info *QosPolicyRenameInfo) (*QosPolicyResponse, *http.Response, error) {
	qp.Params.XMLName = xml.Name{Local: "qos-policy-group-rename"}
	qp.Params.QosPolicyRenameInfo = *info

	return qp.doAPICall()
}

// ChangeIops modifies the iops
func (qp QosPolicy) ChangeIops(iops string, qosPolicyName string) (*QosPolicyResponse, *http.Response, error) {
	qp.Params.XMLName = xml.Name{Local: "qos-policy-group-modify"}
	qp.Params.QosPolicyInfo.PolicyGroup = qosPolicyName
	qp.Params.QosPolicyInfo.MaxThroughput = iops

	return qp.doAPICall()
}

// Delete removes qos policy from the cluster, optionally forcing that deletion
func (qp QosPolicy) Delete(name string, force bool) (*QosPolicyResponse, *http.Response, error) {
	qp.Params.XMLName = xml.Name{Local: "qos-policy-group-delete"}
	qp.Params.QosPolicyInfo.PolicyGroup = name
	qp.Params.QosPolicyInfo.Force = force

	return qp.doAPICall()
}

func (qp QosPolicy) doAPICall() (*QosPolicyResponse, *http.Response, error) {
	r := QosPolicyResponse{}
	res, err := qp.get(qp, &r)
	return &r, res, err
}
