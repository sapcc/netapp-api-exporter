package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pepabo/go-netapp/netapp"
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
	SizeTotal                         string
	SizeAvailable                     string
	SizeUsed                          string
	SizeUsedBySnapshots               string
	SizeAvailableForSnapshots         string
	SnapshotReserveSize               string
	PercentageSizeUsed                string
	PercentageCompressionSpaceSaved   string
	PercentageDeduplicationSpaceSaved string
	PercentageTotalSpaceSaved         string
}

func (f *Filer) GetNetappVolume() (volumes []*NetappVolume, err error) {
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
	if err != nil {
		return nil, err
	} else {
		logger.Printf("%s: %d volumes fetched", f.Host, len(vols))
	}

	for _, vol := range vols {
		nv := &NetappVolume{FilerName: f.Name}
		if vol.VolumeIDAttributes != nil {
			nv.Vserver = vol.VolumeIDAttributes.OwningVserverName
			nv.Volume = vol.VolumeIDAttributes.Name
		}
		if vol.VolumeSpaceAttributes != nil {
			nv.Size = vol.VolumeSpaceAttributes.Size
			nv.SizeAvailable = vol.VolumeSpaceAttributes.SizeAvailable
			nv.SizeTotal = vol.VolumeSpaceAttributes.SizeTotal
			nv.SizeUsed = vol.VolumeSpaceAttributes.SizeUsed
			nv.SizeUsedBySnapshots = vol.VolumeSpaceAttributes.SizeUsedBySnapshots
			nv.SizeAvailableForSnapshots = vol.VolumeSpaceAttributes.SizeAvailableForSnapshots
			nv.SnapshotReserveSize = vol.VolumeSpaceAttributes.SnapshotReserveSize
			nv.PercentageSizeUsed = vol.VolumeSpaceAttributes.PercentageSizeUsed
		} else {
			logger.Printf("%s has no VolumeSpaceAttributes", nv.Volume)
		}
		if vol.VolumeSisAttributes != nil {
			nv.PercentageCompressionSpaceSaved = vol.VolumeSisAttributes.PercentageCompressionSpaceSaved
			nv.PercentageDeduplicationSpaceSaved = vol.VolumeSisAttributes.PercentageDeduplicationSpaceSaved
			nv.PercentageTotalSpaceSaved = vol.VolumeSisAttributes.PercentageTotalSpaceSaved
		} else {
			logger.Printf("%s has no VolumeSisAttributes", vol.VolumeIDAttributes.Name)
			logger.Debugf("%+v", vol.VolumeIDAttributes)
		}
		if vol.VolumeIDAttributes.Comment == "" {
			if !strings.Contains(vol.VolumeIDAttributes.Name, "root") &&
				!strings.Contains(vol.VolumeIDAttributes.Name, "vol0") {
				logger.Printf("%s (%s) does not have comment", vol.VolumeIDAttributes.Name, vol.VolumeIDAttributes.OwningVserverName)
			}
		} else {
			// nv.ShareID, nv.ShareName, nv.ProjectID := parseVolumeComment(vol.VolumeIDAttributes.Comment)
			shareID, shareName, projectID, err := parseVolumeComment(vol.VolumeIDAttributes.Comment)
			if err != nil {
				logger.Warn(err)
			} else {
				nv.ShareID = shareID
				nv.ShareName = shareName
				nv.ProjectID = projectID
			}
		}
		volumes = append(volumes, nv)
	}

	return volumes, err

}

func (f *Filer) getVolumeList(opts *netapp.VolumeOptions) (res []netapp.VolumeInfo, err error) {
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
	// r := regexp.MustCompile(`((?P<k1>\w+): (?P<v1>[\w-]+))(, ((?P<k2>\w+): (?P<v2>[\w-]+))(, ((?P<k3>\w+): (?P<v3>[\w-]+)))?)?`)
	// matches := r.FindStringSubmatch(c)

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
