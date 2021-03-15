package netapp

import (
	"encoding/xml"
	"net/http"
)

type vServerExportsRequest struct {
	Base
	Params struct {
		XMLName               xml.Name
		VServerExportRuleInfo `xml:",innerxml"`
	}
}

type ExportRuleQuery struct {
	ExportRuleInfo *VServerExportRuleInfo `xml:"export-rule-info,omitempty"`
}

// VServerExportRuleInfo sets all different options for Export Rules
type VServerExportRuleInfo struct {
	AnonymousUserID           int       `xml:"anonymous-user-id,omitempty"`
	ClientMatch               string    `xml:"client-match,omitempty"`
	ExportChownMode           string    `xml:"export-chown-mode,omitempty"`
	ExportNTFSUnixSecurityOps string    `xml:"export-ntfs-unix-security-ops,omitempty"`
	AllowCreateDevices        bool      `xml:"is-allow-dev-is-enabled,omitempty"`
	AllowSetUID               bool      `xml:"is-allow-set-uid-enabled,omitempty"`
	PolicyName                string    `xml:"policy-name,omitempty"`
	Protocol                  *[]string `xml:"protocol>access-protocol,omitempty"`
	ReadOnlyRule              *[]string `xml:"ro-rule>security-flavor,omitempty"`
	RuleIndex                 int       `xml:"rule-index,omitempty"`
	ReadWriteRule             *[]string `xml:"rw-rule>security-flavor,omitempty"`
	SuperUserSecurity         *[]string `xml:"super-user-security>security-flavor,omitempty"`
}

// VServerExportsResponse creates correct response obj
type VServerExportsResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
	} `xml:"results"`
}

type ExportRuleListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			VServerExportRuleInfo []VServerExportRuleInfo `xml:"export-rule-info"`
		} `xml:"attributes-list"`
	} `xml:"results"`
}

// ListExportRules lists all export rules on a vserver
func (v VServer) ListExportRules(vServerName string) (*ExportRuleListResponse, *http.Response, error) {
	req := v.newVServerExportsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "export-rule-get-iter"}
	req.Params.VServerExportRuleInfo = VServerExportRuleInfo{}
	r := &ExportRuleListResponse{}
	res, err := req.get(req, r)
	return r, res, err
}

// CreateExportRule creates a new export rule for a given vserver
func (v VServer) CreateExportRule(vServerName string, options *VServerExportRuleInfo) (*VServerExportsResponse, *http.Response, error) {
	req := v.newVServerExportsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "export-rule-create"}
	req.Params.VServerExportRuleInfo = *options

	r := &VServerExportsResponse{}
	res, err := v.get(req, r)
	return r, res, err
}

// DeleteExportRule removes an export rule for a given vserver, policy and rule index
func (v VServer) DeleteExportRule(vServerName string, policyName string, ruleIndex int) (*VServerExportsResponse, *http.Response, error) {
	req := v.newVServerExportsRequest()
	req.Base.Name = vServerName
	req.Params.XMLName = xml.Name{Local: "export-rule-destroy"}
	req.Params.VServerExportRuleInfo = VServerExportRuleInfo{
		PolicyName: policyName,
		RuleIndex:  ruleIndex,
	}

	r := &VServerExportsResponse{}
	res, err := v.get(req, r)
	return r, res, err
}

func (v VServer) newVServerExportsRequest() *vServerExportsRequest {
	return &vServerExportsRequest{
		Base: v.Base,
	}
}
