package netapp

import (
	"encoding/xml"
	"net/http"
)

type Aggregate struct {
	Base
	Params struct {
		XMLName xml.Name
		AggrOptions
	}
}

type AggrOptions struct {
	DesiredAttributes *AggrInfo `xml:"desired-attributes>aggr-attributes,omitempty"`
	MaxRecords        int       `xml:"max-records,omitempty"`
	Query             *AggrInfo `xml:"query>aggr-attributes,omitempty"`
	Tag               string    `xml:"tag,omitempty"`
}

type AggrListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AggrAttributes []AggrInfo `xml:"attributes-list>aggr-attributes"`
		NextTag        string     `xml:"next-tag"`
	} `xml:"results"`
}

func (a Aggregate) List(options *AggrOptions) (*AggrListResponse, *http.Response, error) {
	a.Params.XMLName = xml.Name{Local: "aggr-get-iter"}
	a.Params.AggrOptions = *options
	r := AggrListResponse{}
	res, err := a.get(a, &r)
	return &r, res, err
}

type AggrListPagesResponse struct {
	Response    *AggrListResponse
	Error       error
	RawResponse *http.Response
}

type AggregatePageHandler func(AggrListPagesResponse) (shouldContinue bool)

func (a *Aggregate) ListPages(options *AggrOptions, fn AggregatePageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		aggregateResponse, res, err := a.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(AggrListPagesResponse{Response: aggregateResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = aggregateResponse.Results.NextTag
			requestOptions = &AggrOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}

}

type AggrInfo struct {
	AggregateName           string                   `xml:"aggregate-name,omitempty"`
	AggrInodeAttributes     *AggrInodeAttributes     `xml:"aggr-inode-attributes,omitempty"`
	AggrSpaceAttributes     *AggrSpaceAttributes     `xml:"aggr-space-attributes,omitempty"`
	AggrOwnershipAttributes *AggrOwnershipAttributes `xml:"aggr-ownership-attributes,omitempty"`
	AggrRaidAttributes      *AggrRaidAttributes      `xml:"aggr-raid-attributes,omitempty"`
}

type AggrRaidAttributes struct {
	AggregateType      string `xml:"aggregate-type,omitempty"`
	CacheRaidGroupSize int    `xml:"cache-raid-group-size,omitempty"`
	ChecksumStatus     string `xml:"checksum-status,omitempty"`
	ChecksumStyle      string `xml:"checksum-style,omitempty"`
	DiskCount          int    `xml:"disk-count,omitempty"`
	EncryptionKeyID    string `xml:"encryption-key-id,omitempty"`
	HaPolicy           string `xml:"ha-policy,omitempty"`
	HasLocalRoot       *bool  `xml:"has-local-root"`
	HasPartnerRoot     *bool  `xml:"has-partner-root"`
	IsChecksumEnabled  *bool  `xml:"is-checksum-enabled"`
	IsEncrypted        *bool  `xml:"is-encrypted"`
	IsHybrid           *bool  `xml:"is-hybrid"`
	IsHybridEnabled    *bool  `xml:"is-hybrid-enabled"`
	IsInconsistent     *bool  `xml:"is-inconsistent"`
	IsMirrored         *bool  `xml:"is-mirrored"`
	IsRootAggregate    *bool  `xml:"is-root-aggregate"`
	MirrorStatus       string `xml:"mirror-status,omitempty"`
	MountState         string `xml:"mount-state,omitempty"`
	PlexCount          int    `xml:"plex-count,omitempty"`
	RaidLostWriteState string `xml:"raid-lost-write-state,omitempty"`
	RaidSize           int    `xml:"raid-size,omitempty"`
	RaidStatus         string `xml:"raid-status,omitempty"`
	RaidType           string `xml:"raid-type,omitempty"`
	State              string `xml:"state,omitempty"`
	UsesSharedDisks    *bool  `xml:"uses-shared-disks"`
}

// AggrOwnershipAttributes describe aggregate's ownership
type AggrOwnershipAttributes struct {
	Cluster   string `xml:"cluster"`
	HomeID    int    `xml:"home-id"`
	HomeName  string `xml:"home-name"`
	OwnerID   int    `xml:"owner-id"`
	OwnerName string `xml:"owner-name"`
}

type AggrInodeAttributes struct {
	FilesPrivateUsed         int `xml:"files-private-used"`
	FilesTotal               int `xml:"files-total"`
	FilesUsed                int `xml:"files-used"`
	InodefilePrivateCapacity int `xml:"inodefile-private-capacity"`
	InodefilePublicCapacity  int `xml:"inodefile-public-capacity"`
	MaxfilesAvailable        int `xml:"maxfiles-available"`
	MaxfilesPossible         int `xml:"maxfiles-possible"`
	MaxfilesUsed             int `xml:"maxfiles-used"`
	PercentInodeUsedCapacity int `xml:"percent-inode-used-capacity"`
}

type AggrSpaceAttributes struct {
	AggregateMetadata            string `xml:"aggregate-metadata"`
	HybridCacheSizeTotal         string `xml:"hybrid-cache-size-total"`
	PercentUsedCapacity          string `xml:"percent-used-capacity"`
	PhysicalUsed                 int    `xml:"physical-used"`
	PhysicalUsedPercent          int    `xml:"physical-used-percent"`
	SizeAvailable                int    `xml:"size-available"`
	SizeTotal                    int    `xml:"size-total"`
	SizeUsed                     int    `xml:"size-used"`
	TotalReservedSpace           int    `xml:"total-reserved-space"`
	UsedIncludingSnapshotReserve string `xml:"used-including-snapshot-reserve"`
	VolumeFootprints             string `xml:"volume-footprints"`
}

type AggregateSpace struct {
	Base
	Params struct {
		XMLName xml.Name
		*AggrSpaceOptions
	}
}
type AggrSpaceInfoQuery struct {
	AggrSpaceInfo *AggrSpaceInfo `xml:"space-information,omitempty"`
}

type AggrSpaceOptions struct {
	DesiredAttributes *AggrSpaceInfoQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                 `xml:"max-records,omitempty"`
	Query             *AggrSpaceInfoQuery `xml:"query,omitempty"`
	Tag               string              `xml:"tag,omitempty"`
}

type AggrSpaceInfo struct {
	Aggregate                           string `xml:"aggregate,omitempty"`
	AggregateMetadata                   string `xml:"aggregate-metadata,omitempty"`
	AggregateMetadataPercent            string `xml:"aggregate-metadata-percent,omitempty"`
	AggregateSize                       string `xml:"aggregate-size,omitempty"`
	PercentSnapshotSpace                string `xml:"percent-snapshot-space,omitempty"`
	PhysicalUsed                        string `xml:"physical-used,omitempty"`
	PhysicalUsedPercent                 string `xml:"physical-used-percent,omitempty"`
	SnapSizeTotal                       string `xml:"snap-size-total,omitempty"`
	SnapshotReserveUnusable             string `xml:"snapshot-reserve-unusable,omitempty"`
	SnapshotReserveUnusablePercent      string `xml:"snapshot-reserve-unusable-percent,omitempty"`
	UsedIncludingSnapshotReserve        string `xml:"used-including-snapshot-reserve,omitempty"`
	UsedIncludingSnapshotReservePercent string `xml:"used-including-snapshot-reserve-percent,omitempty"`
	VolumeFootprints                    string `xml:"volume-footprints,omitempty"`
	VolumeFootprintsPercent             string `xml:"volume-footprints-percent,omitempty"`
}

type AggrSpaceListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			AggrAttributes []AggrSpaceInfo `xml:"space-information"`
		} `xml:"attributes-list"`
	} `xml:"results"`
}

func (a *AggregateSpace) List(options *AggrSpaceOptions) (*AggrSpaceListResponse, *http.Response, error) {
	a.Params.XMLName = xml.Name{Local: "aggr-space-get-iter"}
	a.Params.AggrSpaceOptions = options
	r := AggrSpaceListResponse{}
	res, err := a.get(a, &r)
	return &r, res, err
}

type AggregateSpares struct {
	Base
	Params struct {
		XMLName xml.Name
		*AggrSparesOptions
	}
}

func (a *AggregateSpares) List(options *AggrSparesOptions) (*AggrSparesListResponse, *http.Response, error) {
	a.Params.XMLName = xml.Name{Local: "aggr-spare-get-iter"}
	a.Params.AggrSparesOptions = options
	r := AggrSparesListResponse{}
	res, err := a.get(a, &r)
	return &r, res, err
}

func (a *AggregateSpares) ListPages(options *AggrSparesOptions, fn AggregateSparesPageHandler) {

	requestOptions := options

	for shouldContinue := true; shouldContinue; {
		aggregateResponse, res, err := a.List(requestOptions)
		handlerResponse := false

		handlerResponse = fn(AggrSparesListPagesResponse{Response: aggregateResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = aggregateResponse.Results.NextTag
			requestOptions = &AggrSparesOptions{
				Tag:        nextTag,
				MaxRecords: options.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}

type AggrSpareDiskInfoQuery struct {
	AggrSpareDiskInfo *AggrSpareDiskInfo `xml:"aggr-spare-disk-info,omitempty"`
}

type AggrSparesOptions struct {
	DesiredAttributes *AggrSpareDiskInfoQuery `xml:"desired-attributes,omitempty"`
	MaxRecords        int                     `xml:"max-records,omitempty"`
	Query             *AggrSpareDiskInfoQuery `xml:"query,omitempty"`
	Tag               string                  `xml:"tag,omitempty"`
}

type AggregateSparesPageHandler func(AggrSparesListPagesResponse) (shouldContinue bool)

type AggrSpareDiskInfo struct {
	ChecksumStyle           string `xml:"checksum-style"`
	Disk                    string `xml:"disk"`
	DiskRpm                 int    `xml:"disk-rpm"`
	DiskType                string `xml:"disk-type"`
	EffectiveDiskRpm        int    `xml:"effective-disk-rpm"`
	EffectiveDiskType       string `xml:"effective-disk-type"`
	IsDiskLeftBehind        bool   `xml:"is-disk-left-behind"`
	IsDiskShared            bool   `xml:"is-disk-shared"`
	IsDiskZeroed            bool   `xml:"is-disk-zeroed"`
	IsDiskZeroing           bool   `xml:"is-disk-zeroing"`
	IsSparecore             bool   `xml:"is-sparecore"`
	LocalUsableDataSize     int    `xml:"local-usable-data-size"`
	LocalUsableDataSizeBlks int    `xml:"local-usable-data-size-blks"`
	LocalUsableRootSize     int    `xml:"local-usable-root-size"`
	LocalUsableRootSizeBlks int    `xml:"local-usable-root-size-blks"`
	OriginalOwner           string `xml:"original-owner"`
	SyncmirrorPool          string `xml:"syncmirror-pool"`
	TotalSize               int    `xml:"total-size"`
	UsableSize              int    `xml:"usable-size"`
	UsableSizeBlks          int    `xml:"usable-size-blks"`
	ZeroingPercent          int    `xml:"zeroing-percent"`
}

type AggrSparesListResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			AggrAttributes []AggrSpareDiskInfo `xml:"aggr-spare-disk-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type AggrSparesListPagesResponse struct {
	Response    *AggrSparesListResponse
	Error       error
	RawResponse *http.Response
}
