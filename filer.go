package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/pepabo/go-netapp/netapp"
)

const (
	_url    = "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	version = "1.7"
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
	ShareID    string
	ShareName  string
	ProjectId  string
	InstanceID string
}

type NetappVolume struct {
	Vserver            string
	Volume             string
	SizeTotal          float64
	SizeAvailable      float64
	SizeUsed           float64
	PercentageSizeUsed float64
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
	url := fmt.Sprintf(_url, f.Host)
	netappOpt := &netapp.ClientOptions{
		BasicAuthUser:     f.Username,
		BasicAuthPassword: f.Password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}
	f.NetappClient = netapp.NewClient(url, version, netappOpt)

	IDEndpoint := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", os.Getenv("OS_REGION"))

	osOpt := gophercloud.AuthOptions{
		IdentityEndpoint: IDEndpoint,
		DomainName:       "ccadmin",
		TenantName:       "cloud_admin",
		Username:         os.Getenv("OS_USER"),
		Password:         os.Getenv("OS_PASSWORD"),
	}

	provider, err := openstack.AuthenticatedClient(osOpt)
	if err != nil {
		log.Fatal(err)
	}

	eo := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION")}
	f.OpenstackClient, err = openstack.NewSharedFileSystemV2(provider, eo)
}

func (f *Filer) GetManilaShare() map[string]ManilaShare {
	lo := shares.ListOpts{AllTenants: true}
	allpages, err := shares.ListDetail(f.OpenstackClient, lo).AllPages()
	if err != nil {
		log.Fatal(err)
	}

	sh, err := shares.ExtractShares(allpages)
	if err != nil {
		log.Fatal(err)
	}

	r := make(map[string]ManilaShare)

	for _, s := range sh {
		// r[instance_id]
		r[s.Name] = ManilaShare{
			ShareID:   s.ID,
			ShareName: s.Name,
			ProjectId: s.ProjectID,
		}
	}

	return r
}

func (f *Filer) GetNetappVolume() (r []*NetappVolume, err error) {

	vserverOptions := netapp.VServerOptions{
		Query: &netapp.VServerQuery{
			VServerInfo: &netapp.VServerInfo{
				VserverType: "cluster | data",
			},
		},
		DesiredAttributes: &netapp.VServerQuery{
			VServerInfo: &netapp.VServerInfo{
				VserverName: "x",
				UUID:        "x",
			},
		},
		MaxRecords: 100,
	}

	volumeOptions := netapp.VolumeOptions{
		MaxRecords: 200,
		Query: &netapp.VolumeQuery{
			VolumeInfo: &netapp.VolumeInfo{
				VolumeIDAttributes: &netapp.VolumeIDAttributes{
					OwningVserverUUID: "x",
				},
			},
		},
		DesiredAttributes: &netapp.VolumeQuery{
			VolumeInfo: &netapp.VolumeInfo{
				VolumeIDAttributes: &netapp.VolumeIDAttributes{
					Name:              "x",
					OwningVserverName: "x",
					OwningVserverUUID: "x",
				},
				VolumeSpaceAttributes: &netapp.VolumeSpaceAttributes{
					//
					Size:                1,
					SizeTotal:           "x",
					SizeAvailable:       "x",
					SizeUsed:            "x",
					SizeUsedBySnapshots: "x",
					PercentageSizeUsed:  "x",
				},
			},
		},
	}

	vserverList, _, _ := f.NetappClient.VServer.List(&vserverOptions)
	// fmt.Println("vserverList ", vserverList)

	for _, vserver := range vserverList.Results.AttributesList.VserverInfo {
		volumeOptions.Query.VolumeInfo.VolumeIDAttributes.OwningVserverUUID = vserver.UUID
		vols, _, _ := f.NetappClient.Volume.List(&volumeOptions)

		for _, vol := range vols.Results.AttributesList {
			nv := &NetappVolume{Vserver: vserver.VserverName}
			nv.Volume = vol.VolumeIDAttributes.Name
			nv.SizeAvailable, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeAvailable, 64)
			nv.SizeTotal, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeTotal, 64)
			nv.SizeUsed, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.SizeUsed, 64)
			nv.PercentageSizeUsed, err = strconv.ParseFloat(vol.VolumeSpaceAttributes.PercentageSizeUsed, 64)

			r = append(r, nv)
		}
	}

	return
}
