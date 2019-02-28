package main

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func TestC(t *testing.T) {

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
	assert.Equal(t, "", client)

	lo := shares.ListOpts{AllTenants: true}
	allpages, _ := shares.ListDetail(client, lo).AllPages()
	sh, _ := shares.ExtractShares(allpages)

	s := sh[0]
	assert.Equal(t, "1", s.ProjectID)
	assert.Equal(t, "1", s.ShareServerID)
	assert.Equal(t, "1", s.Name)
	assert.Equal(t, "1", s.DisplayName)
	assert.Equal(t, "1", s.ShareTypeName)
}
