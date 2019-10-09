package main

import (
	"fmt"
	"github.com/pepabo/go-netapp/netapp"
	"regexp"
	"strconv"
)

type NetappVolume struct {
	ProjectID                         string
	ShareID                           string
	ShareName                         string
	FilerName                         string
	Vserver                           string
	Volume                            string
	Comment                           string
	Size                              int
	SizeTotal                         float64
	SizeAvailable                     float64
	SizeUsed                          float64
	SizeUsedBySnapshots               float64
	SizeAvailableForSnapshots         float64
	SnapshotReserveSize               float64
	PercentageSizeUsed                float64
	PercentageCompressionSpaceSaved   float64
	PercentageDeduplicationSpaceSaved float64
	PercentageTotalSpaceSaved         float64
}

// GetNetappVolume() returns list of volumes from netapp filer.
func (f *FilerManager) GetNetappVolume() (volumes []*NetappVolume, err error) {
	volumeOptions := netapp.VolumeOptions{
		MaxRecords: 20,
		DesiredAttributes: &netapp.VolumeQuery{
			VolumeInfo: &netapp.VolumeInfo{
				VolumeIDAttributes: &netapp.VolumeIDAttributes{
					Name:              "x",
					OwningVserverName: "x",
					OwningVserverUUID: "x",
					Comment:           "x",
				},
				VolumeSpaceAttributes: &netapp.VolumeSpaceAttributes{
					Size:                      1,
					SizeTotal:                 "x",
					SizeAvailable:             "x",
					SizeUsed:                  "x",
					SizeUsedBySnapshots:       "x",
					SizeAvailableForSnapshots: "x",
					SnapshotReserveSize:       "x",
					PercentageSizeUsed:        "x",
				},
				VolumeSisAttributes: &netapp.VolumeSisAttributes{
					PercentageCompressionSpaceSaved:   "x",
					PercentageDeduplicationSpaceSaved: "x",
					PercentageTotalSpaceSaved:         "x",
				},
			},
		},
	}

	vols, err := f.getVolumeList(&volumeOptions)

	if err == nil {
		logger.Printf("%s: %d volumes fetched", f.Host, len(vols))
		for _, vol := range vols {
			nv := &NetappVolume{FilerName: f.Name}
			if vol.VolumeIDAttributes != nil {
				nv.Vserver = vol.VolumeIDAttributes.OwningVserverName
				nv.Volume = vol.VolumeIDAttributes.Name
			}
			if vol.VolumeSpaceAttributes != nil {
				v := vol.VolumeSpaceAttributes
				sizeTotal, _ := strconv.ParseFloat(v.SizeTotal, 64)
				sizeAvailable, _ := strconv.ParseFloat(v.SizeAvailable, 64)
				sizeUsed, _ := strconv.ParseFloat(v.SizeUsed, 64)
				sizeUsedBySnapshots, _ := strconv.ParseFloat(v.SizeUsedBySnapshots, 64)
				sizeAvailableForSnapshots, _ := strconv.ParseFloat(v.SizeAvailableForSnapshots, 64)
				snapshotReserveSize, _ := strconv.ParseFloat(v.SnapshotReserveSize, 64)
				percentageSizeUsed, _ := strconv.ParseFloat(v.PercentageSizeUsed, 64)

				nv.Size = vol.VolumeSpaceAttributes.Size
				nv.SizeAvailable = sizeAvailable
				nv.SizeTotal = sizeTotal
				nv.SizeUsed = sizeUsed
				nv.SizeUsedBySnapshots = sizeUsedBySnapshots
				nv.SizeAvailableForSnapshots = sizeAvailableForSnapshots
				nv.SnapshotReserveSize = snapshotReserveSize
				nv.PercentageSizeUsed = percentageSizeUsed
			} else {
				logger.Warnf("%s has no VolumeSpaceAttributes", nv.Volume)
			}
			if vol.VolumeSisAttributes != nil {
				v := vol.VolumeSisAttributes
				percentageCompressionSpaceSaved, _ := strconv.ParseFloat(v.PercentageCompressionSpaceSaved, 64)
				percentageDeduplicationSpaceSaved, _ := strconv.ParseFloat(v.PercentageDeduplicationSpaceSaved, 64)
				percentageTotalSpaceSaved, _ := strconv.ParseFloat(v.PercentageTotalSpaceSaved, 64)

				nv.PercentageCompressionSpaceSaved = percentageCompressionSpaceSaved
				nv.PercentageDeduplicationSpaceSaved = percentageDeduplicationSpaceSaved
				nv.PercentageTotalSpaceSaved = percentageTotalSpaceSaved
			} else {
				logger.Warnf("%s has no VolumeSisAttributes", vol.VolumeIDAttributes.Name)
				logger.Debugf("%+v", vol.VolumeIDAttributes)
			}
			if vol.VolumeIDAttributes.Comment != "" {
				shareID, shareName, projectID, err := parseVolumeComment(vol.VolumeIDAttributes.Comment)
				if err != nil {
					logger.Warn(err)
				} else {
					nv.ShareID = shareID
					nv.ShareName = shareName
					nv.ProjectID = projectID
				}
			} else {
				//logger.Warnf("%s (%s) does not have comment",
				//	vol.VolumeIDAttributes.Name, vol.VolumeIDAttributes.OwningVserverName)
			}
			volumes = append(volumes, nv)
		}
	}

	return
}

func (f *FilerManager) getVolumeList(opts *netapp.VolumeOptions) (res []netapp.VolumeInfo, err error) {
	pageHandler := func(r netapp.VolumeListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AttributesList...)
		return true
	}
	f.NetappClient.Volume.ListPages(opts, pageHandler)
	return
}

func parseVolumeComment(c string) (shareID string, shareName string, projectID string, err error) {
	r := regexp.MustCompile(`(\w+): ([\w-]+)`)
	matches := r.FindAllStringSubmatch(c, 3)
	for _, m := range matches {
		switch m[1] {
		case "share_id":
			shareID = m[2]
		case "share_name":
			shareName = m[2]
		case "project":
			projectID = m[2]
		}
	}
	if shareID == "" || projectID == "" {
		err = fmt.Errorf("failed to parse share_id/project from '%s'", c)
	}
	return
}
