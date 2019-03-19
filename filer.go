package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/pepabo/go-netapp/netapp"
)

type Filer struct {
	FilerBase
	NetappClient    *netapp.Client
	OpenstackClient *gophercloud.ServiceClient
}

type FilerBase struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type NetappVolume struct {
	ShareID                           string
	ShareName                         string
	ProjectID                         string
	Vserver                           string
	Volume                            string
	SizeTotal                         string
	SizeAvailable                     string
	SizeUsed                          string
	PercentageSizeUsed                string
	PercentageCompressionSpaceSaved   string
	PercentageDeduplicationSpaceSaved string
	PercentageTotalSpaceSaved         string
}

func NewFiler(name, host, username, password string) *Filer {
	f := &Filer{
		FilerBase: FilerBase{
			Name:     name,
			Host:     host,
			Username: username,
			Password: password,
		},
	}
	f.Init()
	return f
}

func (f *Filer) Init() {
	f.NetappClient = newNetappClient(f.Host, f.Username, f.Password)
}

func newNetappClient(host, username, password string) *netapp.Client {
	_url := "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	url := fmt.Sprintf(_url, host)

	version := "1.7"

	opts := &netapp.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}

	return netapp.NewClient(url, version, opts)
}

func (f *Filer) GetNetappVolume() (r []*NetappVolume, err error) {
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
					Size:                1,
					SizeTotal:           "x",
					SizeAvailable:       "x",
					SizeUsed:            "x",
					SizeUsedBySnapshots: "x",
					PercentageSizeUsed:  "x",
				},
				VolumeSisAttributes: &netapp.VolumeSisAttributes{
					PercentageCompressionSpaceSaved:   "x",
					PercentageDeduplicationSpaceSaved: "x",
					PercentageTotalSpaceSaved:         "x",
				},
			},
		},
	}

	volumePages := f.getNetappVolumePages(&volumeOptions, -1)
	volumes := extracVolumes(volumePages)

	log.Printf("%s: %d (%d) volumes fetched", f.Host, len(volumes), len(volumePages))
	// if len(volumes) > 0 {
	// 	log.Printf("%+v", volumes[0].VolumeIDAttributes)
	// 	log.Printf("%+v", volumes[0].VolumeSpaceAttributes)
	// }

	for _, vol := range volumes {
		nv := &NetappVolume{
			Vserver:            vol.VolumeIDAttributes.OwningVserverName,
			Volume:             vol.VolumeIDAttributes.Name,
			SizeAvailable:      vol.VolumeSpaceAttributes.SizeAvailable,
			SizeTotal:          vol.VolumeSpaceAttributes.SizeTotal,
			SizeUsed:           vol.VolumeSpaceAttributes.SizeUsed,
			PercentageSizeUsed: vol.VolumeSpaceAttributes.PercentageSizeUsed,
		}
		if vol.VolumeSisAttributes != nil {
			nv.PercentageCompressionSpaceSaved = vol.VolumeSisAttributes.PercentageCompressionSpaceSaved
			nv.PercentageDeduplicationSpaceSaved = vol.VolumeSisAttributes.PercentageDeduplicationSpaceSaved
			nv.PercentageTotalSpaceSaved = vol.VolumeSisAttributes.PercentageTotalSpaceSaved
		} else {
			log.Printf("%s has no VolumeSisAttributes", vol.VolumeIDAttributes.Name)
			log.Debugf("%+v", vol.VolumeIDAttributes)
		}

		if vol.VolumeIDAttributes.Comment == "" {
			if vol.VolumeIDAttributes.Name != "root" {
				log.Printf("%s (%s) does not have comment", vol.VolumeIDAttributes.Name, vol.VolumeIDAttributes.OwningVserverName)
			}
		} else {
			nv.ShareID, nv.ShareName, nv.ProjectID = parseComment(vol.VolumeIDAttributes.Comment)
		}

		r = append(r, nv)
	}

	return
}

func parseComment(c string) (shareID string, shareName string, projectID string) {
	// r := regexp.MustCompile(`(share_id:[[:space:]](?P<id>[\-0-9a-z]+))?.*(share_name: (?P<name>[0-9a-zA-Z_]+))?.*(project:[[:space:]](?P<project>[0-9a-z]+))?`)
	// r := regexp.MustCompile(`(share_id: (?P<id>[\-0-9a-z]+))?.*share_name: (?P<name>[0-9A-Za-z_\-]+).* project: (?P<project>\w+)`)
	r := regexp.MustCompile(`((?P<k1>\w+): (?P<v1>[\w-]+))(, ((?P<k2>\w+): (?P<v2>[\w-]+))(, ((?P<k3>\w+): (?P<v3>[\w-]+)))?)?`)

	matches := r.FindStringSubmatch(c)

	for i, m := range matches {
		switch m {
		case "share_id":
			shareID = matches[i+1]
		case "share_name":
			shareName = matches[i+1]
		case "project":
			projectID = matches[i+1]
		}
	}

	return
}

func (f *Filer) getNetappVolumePages(opts *netapp.VolumeOptions, maxPage int) []*netapp.VolumeListResponse {
	var volumePages []*netapp.VolumeListResponse
	var page int

	pageHandler := func(r netapp.VolumeListPagesResponse) bool {
		if r.Error != nil {
			log.Printf("%s", r.Error)
			return false
		}

		volumePages = append(volumePages, r.Response)

		page += 1
		if maxPage > 0 && page >= maxPage {
			return false
		}
		return true
	}

	f.NetappClient.Volume.ListPages(opts, pageHandler)
	return volumePages
}

func extracVolumes(pages []*netapp.VolumeListResponse) (vols []netapp.VolumeInfo) {
	for _, p := range pages {
		vols = append(vols, p.Results.AttributesList...)
	}
	return
}
