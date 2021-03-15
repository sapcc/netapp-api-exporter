package netapp

import (
	"encoding/xml"
	"net/http"
)

// Snapmirror is Snapmirror API struct
type Snapmirror struct {
	Base
	Params struct {
		XMLName           xml.Name
		DesiredAttributes *SnapmirrorInfo `xml:"desired-attributes>snapmirror-info,omitempty"`
		*SnapmirrorInfo
	}
}

type snapmirrorIterRequest struct {
	Base
	Params struct {
		XMLName           xml.Name
		ContinueOnFailure bool            `xml:"continue-on-failure,omitempty"`
		MaxFailureCount   int             `xml:"max-failure-count,omitempty"`
		MaxRecords        int             `xml:"max-records,omitempty"`
		Tag               string          `xml:"tag,omitempty"`
		Query             *SnapmirrorInfo `xml:"query>snapmirror-info"`
	}
}

// Snapmirror Relationship Types
const (
	SnapmirrorRelationshipDP  string = "data_protection"
	SnapmirrorRelationshipLS  string = "load_sharing"
	SnapmirrorRelationshipV   string = "vault"
	SnapmirrorRelationshipR   string = "restore"
	SnapmirrorRelationshipTDP string = "transition_data_protection"
	SnapmirrorRelationshipEDP string = "extended_data_protection"
)

// SnapmirrorInfo contains all fields for snapmirror data
type SnapmirrorInfo struct {
	BreakFailedCount                    int      `xml:"break-failed-count,omitempty"`
	BreakSuccessCount                   int      `xml:"break-successful-count,omitempty"`
	CGItemMappings                      []string `xml:"cg-item-mappings,omitempty"`
	CurrentMaxTransferRate              int      `xml:"current-max-transfer-rate,omitempty"`
	CurrentOperationID                  string   `xml:"current-operation-id,omitempty"`
	CurrentTransferError                string   `xml:"current-transfer-error,omitempty"`
	CurrentTransferPriority             string   `xml:"current-transfer-priority,omitempty"`
	CurrentTransferType                 string   `xml:"current-transfer-type,omitempty"`
	DestinationCluster                  string   `xml:"destination-cluster,omitempty"`
	DestinationLocation                 string   `xml:"destination-location,omitempty"`
	DestinationVolume                   string   `xml:"destination-volume,omitempty"`
	DestinationVolumeNode               string   `xml:"destination-volume-node,omitempty"`
	DestinationVServer                  string   `xml:"destination-vserver,omitempty"`
	ExportedSnapshot                    string   `xml:"exported-snapshot,omitempty"`
	ExportedSnapshotTimestamp           int      `xml:"exported-snapshot-timestamp,omitempty"`
	RestoreFileCount                    int      `xml:"file-restore-file-count,omitempty"`
	RestoreFileList                     []string `xml:"file-restore-file-list,omitempty"`
	IdentitityPreserve                  bool     `xml:"identity-preserve,omitempty"`
	IsConstituent                       bool     `xml:"is-constituent,omitempty"`
	IsHealthy                           bool     `xml:"is-healthy,omitempty"`
	LagTime                             int      `xml:"lag-time,omitempty"`
	LastTransferDuration                int      `xml:"last-transfer-duration,omitempty"`
	LastTransferEndTimestamp            int      `xml:"last-transfer-end-timestamp,omitempty"`
	LastTransferError                   string   `xml:"last-transfer-error,omitempty"`
	LastTransferErrorCodes              []int    `xml:"last-transfer-error-codes,omitempty"`
	LastTransferFrom                    string   `xml:"last-transfer-from,omitempty"`
	LastTransferNetworkCompressionRatio string   `xml:"last-transfer-network-compression-ratio,omitempty"`
	LastTransferSize                    int      `xml:"last-transfer-size,omitempty"`
	LastTransferType                    string   `xml:"last-transfer-type,omitempty"`
	MaxTransferRate                     int      `xml:"max-transfer-rate,omitempty"`
	MirrorState                         string   `xml:"mirror-state,omitempty"`
	NetworkCompressionRatio             string   `xml:"network-compression-ratio,omitempty"`
	NewestSnapshot                      string   `xml:"newest-snapshot,omitempty"`
	NewestSnapshotTimestamp             int      `xml:"newest-snapshot-timestamp,omitempty"`
	Policy                              string   `xml:"policy,omitempty"`
	PolicyType                          string   `xml:"policy-type,omitempty"`
	ProgressLastUpdated                 int      `xml:"progress-last-updated,omitempty"`
	PseudoCommonSnapFailedCount         int      `xml:"pseudo-common-snap-failed-count,omitempty"`
	PseudoCommonSnapSuccessCount        int      `xml:"pseudo-common-snap-success-count,omitempty"`
	RelationshipControlPlane            string   `xml:"relationship-control-plane,omitempty"`
	RelationshipGroupType               string   `xml:"relationship-group-type,omitempty"`
	RelationshipID                      string   `xml:"relationship-id,omitempty"`
	RelationshipProgress                int      `xml:"relationship-progress,omitempty"`
	RelationshipStatus                  string   `xml:"relationship-status,omitempty"`
	RelationshipType                    string   `xml:"relationship-type,omitempty"`
	ResyncAvgTimeSyncCg                 int      `xml:"resync-avg-time-sync-cg,omitempty"`
	ResyncFailedCount                   int      `xml:"resync-failed-count,omitempty"`
	ResyncSuccessCount                  int      `xml:"resync-successful-count,omitempty"`
	Schedule                            string   `xml:"schedule,omitempty"`
	SnapshotCheckpoint                  int      `xml:"snapshot-checkpoint,omitempty"`
	SnapshotProgress                    int      `xml:"snapshot-progress,omitempty"`
	SourceCluster                       string   `xml:"source-cluster,omitempty"`
	SourceLocation                      string   `xml:"source-location,omitempty"`
	SourceVolume                        string   `xml:"source-volume,omitempty"`
	SourceVolumeNode                    string   `xml:"source-volume-node,omitempty"`
	SourceVServer                       string   `xml:"source-vserver,omitempty"`
	TotalTransferBytes                  int      `xml:"total-transfer-bytes,omitempty"`
	TotalTransferTime                   int      `xml:"total-transfer-time,omitempty"`
	TransferSnapshot                    string   `xml:"transfer-snapshot,omitempty"`
	Tries                               string   `xml:"tries,omitempty"`
	UnhealthyReason                     string   `xml:"unhealthy-reason,omitempty"`
	UpdateFailedCount                   int      `xml:"update-failed-count,omitempty"`
	UpdateSuccessCount                  int      `xml:"update-successful-count,omitempty"`
	VServer                             string   `xml:"vserver,omitempty"`
}

// SnapmirrorResponse returns results for snapmirror
type SnapmirrorResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		SingleResultBase
		Info *SnapmirrorInfo `xml:"attributes>snapmirror-info"`
	} `xml:"results"`
}

type SnapmirrorAsyncResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		AsyncResultBase
	} `xml:"results"`
}

type SnapmirrorIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		NumFailed    int `xml:"num-failed"`
		NumSucceeded int `xml:"num-succeeded"`
		FailureList  []struct {
			ErrorNo int             `xml:"error-code"`
			Reason  string          `xml:"error-message"`
			Info    *SnapmirrorInfo `xml:"snapmirror-key>snapmirror-info"`
		} `xml:"failure-list>snapmirror-destroy-iter-info"`
		SuccessList []struct {
			Info *SnapmirrorInfo `xml:"snapmirror-key>snapmirror-info"`
		} `xml:"success-list>snapmirror-destroy-iter-info"`
	} `xml:"results"`
}

// Create creates a snapmirror on a vserver with attributes provided. Note, not all attributes
// are supported, refer to docs or api errors to diagnose
func (s Snapmirror) Create(vServerName string, attributes *SnapmirrorInfo) (*SingleResultResponse, *http.Response, error) {
	s.Name = vServerName
	s.Params.XMLName = xml.Name{Local: "snapmirror-create"}
	s.Params.SnapmirrorInfo = attributes

	r := &SingleResultResponse{}
	res, err := s.get(s, r)
	return r, res, err
}

// Get returns data related to a snapmirror
func (s Snapmirror) Get(vServerName string, sourcePath string, destinationPath string, attributes *SnapmirrorInfo) (*SnapmirrorResponse, *http.Response, error) {
	s.Name = vServerName
	s.Params.XMLName = xml.Name{Local: "snapmirror-get"}
	s.Params.SnapmirrorInfo = &SnapmirrorInfo{
		DestinationLocation: destinationPath,
		SourceLocation:      sourcePath,
	}
	if attributes == nil {
		// base response includes source/destination fields only
		s.Params.DesiredAttributes = &SnapmirrorInfo{
			IsHealthy:        true,
			VServer:          " ",
			RelationshipType: " ",
		}
	} else {
		s.Params.DesiredAttributes = attributes
	}
	r := &SnapmirrorResponse{}
	res, err := s.get(s, r)
	return r, res, err
}

func (s Snapmirror) DestroyBy(query *SnapmirrorInfo, continueOnFailure bool) (*SnapmirrorIterResponse, *http.Response, error) {
	req := &snapmirrorIterRequest{
		Base: s.Base,
	}
	req.Params.XMLName = xml.Name{Local: "snapmirror-destroy-iter"}
	req.Params.Query = query
	req.Params.ContinueOnFailure = continueOnFailure
	req.Params.MaxRecords = 20

	r := &SnapmirrorIterResponse{}
	res, err := s.get(req, r)
	return r, res, err
}

func (s Snapmirror) AbortBy(query *SnapmirrorInfo, continueOnFailure bool) (*SnapmirrorIterResponse, *http.Response, error) {
	req := &snapmirrorIterRequest{
		Base: s.Base,
	}
	req.Params.XMLName = xml.Name{Local: "snapmirror-abort-iter"}
	req.Params.Query = query
	req.Params.ContinueOnFailure = continueOnFailure
	req.Params.MaxRecords = 20

	r := &SnapmirrorIterResponse{}
	res, err := s.get(req, r)
	return r, res, err
}

// InitializeLSSet starts a Load Sharing set via an async call
func (s Snapmirror) InitializeLSSet(vServerName string, sourcePath string) (*SnapmirrorAsyncResponse, *http.Response, error) {
	s.Name = vServerName
	s.Params.XMLName = xml.Name{Local: "snapmirror-initialize-ls-set"}
	s.Params.SnapmirrorInfo = &SnapmirrorInfo{
		SourceLocation: sourcePath,
	}

	r := &SnapmirrorAsyncResponse{}
	res, err := s.get(s, r)
	return r, res, err
}

// UpdateLSSet starts a Load Sharing set via an async call
func (s Snapmirror) UpdateLSSet(vServerName string, sourcePath string) (*SnapmirrorAsyncResponse, *http.Response, error) {
	s.Name = vServerName
	s.Params.XMLName = xml.Name{Local: "snapmirror-update-ls-set"}
	s.Params.SnapmirrorInfo = &SnapmirrorInfo{
		SourceLocation: sourcePath,
	}

	r := &SnapmirrorAsyncResponse{}
	res, err := s.get(s, r)
	return r, res, err
}
