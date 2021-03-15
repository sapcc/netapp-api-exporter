package netapp

import (
	"encoding/xml"
	"net/http"
)

type StorageDisk struct {
	Base
	Params struct {
		XMLName xml.Name
		*StorageDiskOptions
	}
}

type StorageDiskInfo struct {
	DiskInventoryInfo *DiskInventoryInfo `xml:"disk-inventory-info,omitempty"`
	DiskName          string             `xml:"disk-name,omitempty"`
	DiskOwnershipInfo *DiskOwnershipInfo `xml:"disk-ownership-info,omitempty"`
}

type DiskInventoryInfo struct {
	BytesPerSector                 int    `xml:"bytes-per-sector,omitempty"`
	CapacitySectors                int    `xml:"capacity-sectors,omitempty"`
	ChecksumCompatibility          string `xml:"checksum-compatibility,omitempty"`
	DiskClusterName                string `xml:"disk-cluster-name,omitempty"`
	DiskType                       string `xml:"disk-type,omitempty"`
	DiskUid                        string `xml:"disk-uid,omitempty"`
	FirmwareRevision               string `xml:"firmware-revision,omitempty"`
	GrownDefectListCount           int    `xml:"grown-defect-list-count,omitempty"`
	HealthMonitorTimeInterval      int    `xml:"health-monitor-time-interval,omitempty"`
	ImportInProgress               *bool  `xml:"import-in-progress,omitempty"`
	IsDynamicallyQualified         *bool  `xml:"is-dynamically-qualified,omitempty"`
	IsMultidiskCarrier             *bool  `xml:"is-multidisk-carrier,omitempty"`
	IsShared                       *bool  `xml:"is-shared,omitempty"`
	MediaScrubCount                int    `xml:"media-scrub-count,omitempty"`
	MediaScrubLastDoneTimeInterval int    `xml:"media-scrub-last-done-time-interval,omitempty"`
	Model                          string `xml:"model,omitempty"`
	ReservationKey                 string `xml:"reservation-key,omitempty"`
	ReservationType                string `xml:"reservation-type,omitempty"`
	RightSizeSectors               int    `xml:"right-size-sectors,omitempty"`
	Rpm                            int    `xml:"rpm,omitempty"`
	SerialNumber                   string `xml:"serial-number,omitempty"`
	Shelf                          string `xml:"shelf,omitempty"`
	ShelfBay                       string `xml:"shelf-bay,omitempty"`
	ShelfUid                       string `xml:"shelf-uid,omitempty"`
	StackID                        int    `xml:"stack-id,omitempty"`
	Vendor                         string `xml:"vendor,omitempty"`
}

type DiskOwnershipInfo struct {
	DiskUid          string `xml:"disk-uid,omitempty"`
	DrHomeNodeId     int    `xml:"dr-home-node-id,omitempty"`
	DrHomeNodeName   string `xml:"dr-home-node-name,omitempty"`
	HomeNodeId       int    `xml:"home-node-id,omitempty"`
	HomeNodeName     string `xml:"home-node-name,omitempty"`
	IsFailed         *bool  `xml:"is-failed,omitempty"`
	OwnerNodeId      int    `xml:"owner-node-id,omitempty"`
	OwnerNodeName    string `xml:"owner-node-name,omitempty"`
	Pool             int    `xml:"pool,omitempty"`
	ReservedByNodeId int    `xml:"reserved-by-node-id,omitempty"`
}

type StorageDiskGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`
	Results struct {
		ResultBase
		AttributesList struct {
			StorageDiskInfo []StorageDiskInfo `xml:"storage-disk-info"`
		} `xml:"attributes-list"`
		NextTag    string `xml:"next-tag"`
		NumRecords int    `xml:"num-records"`
	} `xml:"results"`
}

type StorageDiskInfoPageResponse struct {
	Response    *StorageDiskGetIterResponse
	Error       error
	RawResponse *http.Response
}

type StorageDiskOptions struct {
	DesiredAttributes *StorageDiskInfo `xml:"desired-attributes>storage-disk-info,omitempty"`
	Query             *StorageDiskInfo `xml:"query>storage-disk-info,omitempty"`
	MaxRecords        int              `xml:"max-records,omitempty"`
	Tag               string           `xml:"tag,omitempty"`
}

func (s *StorageDisk) StorageDiskGetIter(options *StorageDiskOptions) (*StorageDiskGetIterResponse, *http.Response, error) {
	s.Params.XMLName = xml.Name{Local: "storage-disk-get-iter"}
	s.Params.StorageDiskOptions = options
	r := StorageDiskGetIterResponse{}
	res, err := s.get(s, &r)
	return &r, res, err
}

type StorageDiskGetAllPageHandler func(StorageDiskInfoPageResponse) (shouldContinue bool)

func (s *StorageDisk) StorageDiskGetAll(options *StorageDiskOptions, fn StorageDiskGetAllPageHandler) {

	requestOptions := options
	for shouldContinue := true; shouldContinue; {
		storageDiskGetIterResponse, res, err := s.StorageDiskGetIter(requestOptions)
		handlerResponse := false

		handlerResponse = fn(StorageDiskInfoPageResponse{Response: storageDiskGetIterResponse, Error: err, RawResponse: res})

		nextTag := ""
		if err == nil {
			nextTag = storageDiskGetIterResponse.Results.NextTag
			requestOptions = &StorageDiskOptions{
				Tag:        nextTag,
				MaxRecords: requestOptions.MaxRecords,
			}
		}
		shouldContinue = nextTag != "" && handlerResponse
	}
}
