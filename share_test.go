package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func TestC(t *testing.T) {
	var err error

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://identity-3.staging.cloud.sap/v3",
		DomainName:       "ccadmin",
		TenantName:       "cloud_admin",
		Username:         os.Getenv("OS_USER"),
		Password:         os.Getenv("OS_PASSWORD"),
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if assert.Nil(t, err) {

		eo := gophercloud.EndpointOpts{Region: "staging"}
		client, err := openstack.NewSharedFileSystemV2(provider, eo)
		if assert.Nil(t, err) {
			// assert.Equal(t, "", client)

			lo := shares.ListOpts{AllTenants: true}
			allpages, err := shares.ListDetail(client, lo).AllPages()
			if assert.Nil(t, err) {

				sh, _ := shares.ExtractShares(allpages)
				s := sh[0]
				fmt.Println("ID\t\t", s.ID)
				fmt.Println("Name\t\t", s.Name)
				fmt.Println("ProjectID\t", s.ProjectID)
				fmt.Println("ShareServerID\t", s.ShareServerID)
				fmt.Println("ShareType\t", s.ShareType)
				fmt.Println("AvailabilityZone", s.AvailabilityZone)
				fmt.Printf("%T %+v", s, s)
			}
		}
	}

	if err != nil {
		assert.Equal(t, "", err.Error())
	}
}
