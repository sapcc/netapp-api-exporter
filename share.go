package main

import (
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

type ProjectShare struct {
	Project   string
	ShareName string
}

type ProjectShareMap struct {
	Data map[string]ProjectShare
}

func (m *ProjectShareMap) Get() {

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://identity-3.staging.cloud.sap/v3",
		DomainName:       "ccadmin",
		TenantName:       "cloud_admin",
		Username:         os.Getenv("OS_USER"),
		Password:         os.Getenv("OS_PASSWORD"),
	}

	eo := gophercloud.EndpointOpts{Region: "staging"}
	provider, _ := openstack.AuthenticatedClient(opts)
	client, _ := openstack.NewSharedFileSystemV2(provider, eo)

	lo := shares.ListOpts{AllTenants: true}
	allpages, _ := shares.ListDetail(client, lo).AllPages()
	sh, _ := shares.ExtractShares(allpages)

	for _, s := range sh {
		m.Data[s.ProjectID] = ProjectShare{
			Project:   s.ProjectID,
			ShareName: s.Name,
		}
	}
}
