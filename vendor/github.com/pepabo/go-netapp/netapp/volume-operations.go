package netapp

import (
	"encoding/xml"
	"net/http"
)

// These consts are for defined volume operations
const (
	VolumeCreateOperation   = "volume-create"
	VolumeOfflineOperation  = "volume-offline"
	VolumeOnlineOperation   = "volume-online"
	VolumeDestroyOperation  = "volume-destroy"
	VolumeUnmountOperation  = "volume-unmount"
	VolumeRestrictOperation = "volume-restrict"
)

// VolumeOperation is the base struct for volume operations
type VolumeOperation struct {
	Base
	Params struct {
		XMLName    xml.Name
		VolumeName *volumeName
		VolumeCreateOptions
	}
}

type volumeName struct {
	XMLName xml.Name
	Name    string `xml:",innerxml"`
}

// VolumeCreateOptions struct is used for volume creation
type VolumeCreateOptions struct {
	AntivirusOnAccessPolicy    string `xml:"antivirus-on-access-policy,omitempty"`
	CacheRetentionPriority     string `xml:"cache-retention-priority,omitempty"`
	CachingPolicy              string `xml:"caching-policy,omitempty"`
	ConstituentRole            string `xml:"constituent-role,omitempty"`
	ContainingAggregateName    string `xml:"containing-aggr-name,omitempty"`
	EfficiencyPolicy           string `xml:"efficiency-policy,omitempty"`
	Encrypt                    bool   `xml:"encrypt,omitempty"`
	ExcludedFromAutobalance    bool   `xml:"excluded-from-autobalance,omitempty"`
	ExportPolicy               string `xml:"export-policy,omitempty"`
	ExtentSize                 string `xml:"extent-size,omitempty"`
	FlexcachePolicy            string `xml:"flexcache-cache-policy,omitempty"`
	FlexcacheFillPolicy        string `xml:"flexcache-fill-policy,omitempty"`
	FlexcacheOriginVolumeName  string `xml:"flexcache-origin-volume-name,omitempty"`
	GroupID                    int    `xml:"group-id,omitempty"`
	IsJunctionActive           bool   `xml:"is-junction-active,omitempty"`
	IsNvfailEnabled            string `xml:"is-nvfail-enabled,omitempty"`
	IsVserverRoot              bool   `xml:"is-vserver-root,omitempty"`
	JunctionPath               string `xml:"junction-path,omitempty"`
	LanguageCode               string `xml:"language-code,omitempty"`
	MaxDirSize                 int    `xml:"max-dir-size,omitempty"`
	MaxWriteAllocBlocks        int    `xml:"max-write-alloc-blocks,omitempty"`
	PercentageSnapshotReserve  int    `xml:"percentage-snapshot-reserve,omitempty"`
	QosAdaptivePolicyGroupName string `xml:"qos-adaptive-policy-group-name,omitempty"`
	QosPolicyGroupName         string `xml:"qos-policy-group-name,omitempty"`
	Size                       string `xml:"size,omitempty"`
	SnapshotPolicy             string `xml:"snapshot-policy,omitempty"`
	SpaceReserve               string `xml:"space-reserve,omitempty"`
	SpaceSlo                   string `xml:"space-slo,omitempty"`
	StorageService             string `xml:"storage-service,omitempty"`
	TieringPolicy              string `xml:"tiering-policy,omitempty"`
	UnixPermissions            string `xml:"unix-permissions,omitempty"`
	UserID                     int    `xml:"user-id,omitempty"`
	VMAlignSector              int    `xml:"vm-align-sector,omitempty"`
	VMAlignSuffix              string `xml:"vm-align-suffix,omitempty"`
	Volume                     string `xml:"volume,omitempty"`
	VolumeComment              string `xml:"volume-comment,omitempty"`
	VolumeSecurityStyle        string `xml:"volume-security-style,omitempty"`
	VolumeState                string `xml:"volume-state,omitempty"`
	VolumeType                 string `xml:"volume-type,omitempty"`
	VserverDrProtection        string `xml:"vserver-dr-protection,omitempty"`
}

// Create a new volume
func (v VolumeOperation) Create(vserverName string, options *VolumeCreateOptions) (*SingleResultResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: VolumeCreateOperation}
	v.Name = vserverName
	v.Params.VolumeCreateOptions = *options
	r := SingleResultResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}

// Operation runs several operations (from consts defined above with VolumeOperation* name)
func (v VolumeOperation) Operation(vserverName string, volName string, operation string) (*SingleResultResponse, *http.Response, error) {
	v.Params.XMLName = xml.Name{Local: operation}
	v.Name = vserverName
	elementName := "name"
	if operation == VolumeUnmountOperation {
		elementName = "volume-name"
	}
	v.Params.VolumeName = &volumeName{
		XMLName: xml.Name{Local: elementName},
		Name:    volName,
	}
	r := SingleResultResponse{}
	res, err := v.get(v, &r)
	return &r, res, err
}
