package netapp

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	n "github.com/pepabo/go-netapp/netapp"
)

type Volume struct {
	ProjectID                         string
	ShareID                           string
	ShareName                         string
	ShareType                         string
	FilerName                         string
	Vserver                           string
	Volume                            string
	VolumeType                        string
	Comment                           string
	State                             int
	Size                              int
	SizeTotal                         float64
	SizeAvailable                     float64
	SizeUsed                          float64
	SizeUsedBySnapshots               float64
	SizeAvailableForSnapshots         float64
	SnapshotReserveSize               float64
	PercentageSizeUsed                float64
	PercentageSnapshotReserve         float64
	PercentageCompressionSpaceSaved   float64
	PercentageDeduplicationSpaceSaved float64
	PercentageTotalSpaceSaved         float64
	IsEncrypted                       bool
}

func (c *Client) ListVolumes() (volumes []*Volume, err error) {
	volumeInfos, err := c.listVolumes()
	if err != nil {
		return nil, err
	}
	for _, vol := range volumeInfos {
		parsedVol, e := parseVolume(vol)
		if e != nil {
			println(e)
		}
		volumes = append(volumes, parsedVol)
	}
	return
}

func (c *Client) listVolumes() (res []n.VolumeInfo, err error) {
	opts := newVolumeOpts(20)
	pageHandler := func(r n.VolumeListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AttributesList...)
		return true
	}
	c.Volume.ListPages(opts, pageHandler)
	return
}

func newVolumeOpts(maxRecords int) *n.VolumeOptions {
	return &n.VolumeOptions{
		MaxRecords: maxRecords,
		DesiredAttributes: &n.VolumeQuery{
			VolumeInfo: &n.VolumeInfo{
				Encrypt: "x",
				VolumeIDAttributes: &n.VolumeIDAttributes{
					Name:              "x",
					OwningVserverName: "x",
					OwningVserverUUID: "x",
					Comment:           "x",
					Type:              "x",
				},
				VolumeSpaceAttributes: &n.VolumeSpaceAttributes{
					Size:                      1,
					SizeTotal:                 "x",
					SizeAvailable:             "x",
					SizeUsed:                  "x",
					SizeUsedBySnapshots:       "x",
					SizeAvailableForSnapshots: "x",
					SnapshotReserveSize:       "x",
					PercentageSizeUsed:        "x",
					PercentageSnapshotReserve: "x",
				},
				VolumeSisAttributes: &n.VolumeSisAttributes{
					PercentageCompressionSpaceSaved:   "x",
					PercentageDeduplicationSpaceSaved: "x",
					PercentageTotalSpaceSaved:         "x",
				},
				VolumeStateAttributes: &n.VolumeStateAttributes{
					State: "x",
				},
			},
		},
	}
}

func parseVolume(volumeInfo n.VolumeInfo) (*Volume, error) {
	volume := Volume{}
	if volumeInfo.VolumeIDAttributes != nil {
		volume.Vserver = volumeInfo.VolumeIDAttributes.OwningVserverName
		volume.Volume = volumeInfo.VolumeIDAttributes.Name
		volume.VolumeType = volumeInfo.VolumeIDAttributes.Type
	} else {
		msg := fmt.Sprintf("missing VolumeIDAttribtues in %+v", volumeInfo)
		return nil, errors.New(msg)
	}
	if volumeInfo.VolumeSpaceAttributes != nil {
		attributes := volumeInfo.VolumeSpaceAttributes
		sizeTotal, _ := strconv.ParseFloat(attributes.SizeTotal, 64)
		sizeAvailable, _ := strconv.ParseFloat(attributes.SizeAvailable, 64)
		sizeUsed, _ := strconv.ParseFloat(attributes.SizeUsed, 64)
		sizeUsedBySnapshots, _ := strconv.ParseFloat(attributes.SizeUsedBySnapshots, 64)
		sizeAvailableForSnapshots, _ := strconv.ParseFloat(attributes.SizeAvailableForSnapshots, 64)
		snapshotReserveSize, _ := strconv.ParseFloat(attributes.SnapshotReserveSize, 64)
		percentageSizeUsed, _ := strconv.ParseFloat(attributes.PercentageSizeUsed, 64)
		percentageSnapshotReserve, _ := strconv.ParseFloat(attributes.PercentageSnapshotReserve, 64)
		// assign parsed values to output
		volume.Size = attributes.Size
		volume.SizeAvailable = sizeAvailable
		volume.SizeTotal = sizeTotal
		volume.SizeUsed = sizeUsed
		volume.SizeUsedBySnapshots = sizeUsedBySnapshots
		volume.SizeAvailableForSnapshots = sizeAvailableForSnapshots
		volume.SnapshotReserveSize = snapshotReserveSize
		volume.PercentageSizeUsed = percentageSizeUsed
		volume.PercentageSnapshotReserve = percentageSnapshotReserve
	}
	if volumeInfo.VolumeSisAttributes != nil {
		v := volumeInfo.VolumeSisAttributes
		percentageCompressionSpaceSaved, _ := strconv.ParseFloat(v.PercentageCompressionSpaceSaved, 64)
		percentageDeduplicationSpaceSaved, _ := strconv.ParseFloat(v.PercentageDeduplicationSpaceSaved, 64)
		percentageTotalSpaceSaved, _ := strconv.ParseFloat(v.PercentageTotalSpaceSaved, 64)
		// assign parsed values to output
		volume.PercentageCompressionSpaceSaved = percentageCompressionSpaceSaved
		volume.PercentageDeduplicationSpaceSaved = percentageDeduplicationSpaceSaved
		volume.PercentageTotalSpaceSaved = percentageTotalSpaceSaved
	}
	if volumeInfo.VolumeIDAttributes.Comment != "" {
		shareID, shareName, shareType, projectID, err := parseVolumeComment(volumeInfo.VolumeIDAttributes.Comment)
		if err == nil {
			volume.ShareID = shareID
			volume.ShareName = shareName
			volume.ShareType = shareType
			volume.ProjectID = projectID
		}
	}
	if volumeInfo.VolumeStateAttributes != nil {
		if volumeInfo.VolumeStateAttributes.State == "online" {
			volume.State = 1
		} else if volumeInfo.VolumeStateAttributes.State == "restricted" {
			volume.State = 2
		} else if volumeInfo.VolumeStateAttributes.State == "offline" {
			volume.State = 3
		} else if volumeInfo.VolumeStateAttributes.State == "quiesced" {
			volume.State = 4
		}
	}
	if volumeInfo.Encrypt == "true" {
		volume.IsEncrypted = true
	} else {
		volume.IsEncrypted = false
	}
	return &volume, nil
}

func parseVolumeComment(c string) (shareID, shareName, shareType, projectID string, err error) {
	r := regexp.MustCompile(`(\w+): ([\w-]+)`)
	matches := r.FindAllStringSubmatch(c, 4)
	for _, m := range matches {
		switch m[1] {
		case "share_id":
			shareID = m[2]
		case "share_name":
			shareName = m[2]
		case "share_type":
			shareType = m[2]
		case "project":
			projectID = m[2]
		}
	}
	if shareID == "" || projectID == "" {
		err = fmt.Errorf("failed to parse share_id/project from '%s'", c)
	}
	return
}
