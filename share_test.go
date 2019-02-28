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

	IDEndpoint := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", os.Getenv("OS_REGION"))

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: IDEndpoint,
		DomainName:       "ccadmin",
		TenantName:       "cloud_admin",
		Username:         os.Getenv("OS_USER"),
		Password:         os.Getenv("OS_PASSWORD"),
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if assert.Nil(t, err) {

		eo := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION")}
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
				// fmt.Printf("%T %+v", s, s)

				sel, err := shares.GetExportLocations(client, s.ID).Extract()
				fmt.Println(err.Error())
				assert.Nil(t, err)
				assert.Equal(t, "", sel)
				// fmt.Println("ShareInstanceID\t", sel[0].ShareInstanceID)

			}
		}
	}

	if err != nil {
		assert.Equal(t, "", err.Error())
	}
}

// /v2/ae63ddf2076d4342a56eb049e37a7621/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_location
// https://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations
// http://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations
// http://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations
