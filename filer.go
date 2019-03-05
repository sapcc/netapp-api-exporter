package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/pepabo/go-netapp/netapp"

	"github.com/motemen/go-loghttp"
)

type Filer struct {
	FilerBase
	NetappClient    *netapp.Client
	OpenstackClient *gophercloud.ServiceClient
	// Share           *ProjectShareMap
}

type FilerBase struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ManilaShare struct {
	ShareID       string
	ShareName     string
	ShareServerID string
	ProjectId     string
	InstanceID    string
}

type NetappVolume struct {
	Vserver                           string
	Volume                            string
	SizeTotal                         float64
	SizeAvailable                     float64
	SizeUsed                          float64
	PercentageSizeUsed                float64
	PercentageCompressionSpaceSaved   float64
	PercentageDeduplicationSpaceSaved float64
	PercentageTotalSpaceSaved         float64
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

func newManilaClient() (*gophercloud.ServiceClient, error) {
	region := os.Getenv("OS_REGION")
	identityEndpoint := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", region)

	client, err := openstack.NewClient(identityEndpoint)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{}
	config.InsecureSkipVerify = true

	var transport http.RoundTripper
	if os.Getenv("DEBUG") != "" {
		transport = &loghttp.Transport{
			Transport:  &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config},
			LogRequest: logHttpRequestWithHeader,
		}
	} else {
		transport = &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	}

	client.HTTPClient.Transport = transport

	opts := gophercloud.AuthOptions{
		DomainName: "ccadmin",
		TenantName: "cloud_admin",
		Username:   os.Getenv("OS_USERNAME"),
		Password:   os.Getenv("OS_PASSWORD"),
	}

	err = openstack.Authenticate(client, opts)
	if err != nil {
		log.Printf("%+v", opts)
		return nil, err
	}

	eo := gophercloud.EndpointOpts{Region: region}

	manilaClient, err := openstack.NewSharedFileSystemV2(client, eo)
	if err != nil {
		return nil, err
	}

	manilaClient.Microversion = "2.46"
	return manilaClient, nil
}

func (f *Filer) GetManilaShare() (map[string]ManilaShare, error) {
	lo := shares.ListOpts{AllTenants: true}
	allpages, err := shares.ListDetail(f.OpenstackClient, lo).AllPages()
	if err != nil {
		return nil, err
	}

	sh, err := shares.ExtractShares(allpages)
	if err != nil {
		return nil, err
	}

	if os.Getenv("INFO") != "" {
		log.Printf("%d Manila Shares fetched", len(sh))
	}

	r := make(map[string]ManilaShare)
	for _, s := range sh {

		l, err := shares.GetExportLocations(f.OpenstackClient, s.ID).Extract()
		if err != nil {
			return nil, err
		}

		if len(l) > 0 {
			siid := l[0].ShareInstanceID
			siid = strings.Replace(siid, "-", "_", -1)
			r[siid] = ManilaShare{
				ShareID:       s.ID,
				ShareName:     s.Name,
				ShareServerID: s.ShareServerID,
				ProjectId:     s.ProjectID,
				InstanceID:    siid,
			}
		}
	}

	return r, nil
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
			Vserver: vol.VolumeIDAttributes.OwningVserverName,
			Volume:  vol.VolumeIDAttributes.Name,
		}
		nv.SizeAvailable, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeAvailable, 64)
		nv.SizeTotal, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeTotal, 64)
		nv.SizeUsed, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeUsed, 64)
		nv.PercentageSizeUsed, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.PercentageSizeUsed, 64)
		nv.PercentageCompressionSpaceSaved, err = strconv.ParseFloat(vol.VolumeSisAttributes.PercentageCompressionSpaceSaved, 64)
		nv.PercentageDeduplicationSpaceSaved, err = strconv.ParseFloat(vol.VolumeSisAttributes.PercentageDeduplicationSpaceSaved, 64)
		nv.PercentageTotalSpaceSaved, err = strconv.ParseFloat(vol.VolumeSisAttributes.PercentageTotalSpaceSaved, 64)

		r = append(r, nv)
	}

	return
}

func logHttpRequestWithHeader(req *http.Request) {
	log.Printf("--> %s %s %s", req.Method, req.URL, req.Header)
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
