package main

import (
	"fmt"
	"log"
	"os"
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

	manilaClient, err := newManilaClient()
	if err != nil {
		log.Fatal(err)
	}
	f.OpenstackClient = manilaClient
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

	if os.Getenv("INFO") != "" {
		log.Printf("%d volume pages fetched", len(volumePages))
		log.Printf("%d volumes extracted", len(volumes))
		// if len(volumes) > 0 {
		// 	log.Printf("%+v", volumes[0].VolumeIDAttributes)
		// 	log.Printf("%+v", volumes[0].VolumeSpaceAttributes)
		// }
	}

	for _, vol := range volumes {
		nv := &NetappVolume{
			Vserver:                           vol.VolumeIDAttributes.OwningVserverName,
			Volume:                            vol.VolumeIDAttributes.Name,
			SizeAvailable:                     vol.VolumeSpaceAttributes.SizeAvailable,
			SizeTotal:                         vol.VolumeSpaceAttributes.SizeTotal,
			SizeUsed:                          vol.VolumeSpaceAttributes.SizeUsed,
			PercentageSizeUsed:                vol.VolumeSpaceAttributes.PercentageSizeUsed,
			PercentageCompressionSpaceSaved:   vol.VolumeSisAttributes.PercentageCompressionSpaceSaved,
			PercentageDeduplicationSpaceSaved: vol.VolumeSisAttributes.PercentageDeduplicationSpaceSaved,
			PercentageTotalSpaceSaved:         vol.VolumeSisAttributes.PercentageTotalSpaceSaved,
		}

		nv.ShareID, nv.ProjectID = parseComment(vol.VolumeIDAttributes.Comment)

		r = append(r, nv)
	}

	return
}

func parseComment(c string) (shareID string, projectID string) {
	r := regexp.MustCompile(`share_id:[[:space:]](?P<id>[\-0-9a-z]+).*project:[[:space:]]([0-9a-z]+)`)
	matches := r.FindStringSubmatch(c)

	for i, m := range matches {
		switch i {
		case 1:
			shareID = m
		case 2:
			projectID = m
		}
	}

	return
}

func (f *Filer) getNetappVolumePages(opts *netapp.VolumeOptions, maxPage int) []*netapp.VolumeListResponse {
	var volumePages []*netapp.VolumeListResponse
	var page int

	pageHandler := func(r netapp.VolumeListPagesResponse) bool {
		if r.Error != nil {
			if os.Getenv("INFO") != "" {
				log.Printf("%s", r.Error)
			}
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
