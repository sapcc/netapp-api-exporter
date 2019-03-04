package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/pepabo/go-netapp/netapp"

	"github.com/stretchr/testify/assert"
)

func TestNetappVserver(t *testing.T) {
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	c := newNetappClient(host, username, password)

	v := c.VServer

	p := &netapp.VServerOptions{
		DesiredAttributes: &netapp.VServerQuery{
			VServerInfo: &netapp.VServerInfo{
				VserverName:   "1",
				UUID:          "1",
				State:         "1",
				AggregateList: &[]string{"x"},
			},
		},
		MaxRecords: 100,
	}

	r, _, _ := v.List(p)
	// r.Results.NumRecords <= 100

	assert.True(t, r.Results.Passed())
	assert.NotNil(t, r.Results.AttributesList.VserverInfo[0])
}

func TestNetappVolume(t *testing.T) {
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	c := newNetappClient(host, username, password)

	v := c.VolumeSpace
	p := netapp.VolumeSpaceOptions{
		MaxRecords: 10,
	}

	r, _, _ := v.List(&p)

	assert.True(t, r.Results.Passed())
	// assert.NotNil(t, r.Results)
}

func TestManilaClient(t *testing.T) {
	var err error

	// client.Endpoint = "http" + client.Endpoint[5:]
	// a := "http://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations"
	// b := "http://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations"
	// assert.Equal(t, a, b)

	client, err := newManilaClient()
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

	if err != nil {
		assert.Equal(t, "", err.Error())
	}
}

// /v2/ae63ddf2076d4342a56eb049e37a7621/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_location
// https://share-3.staging.cloud.sap/v2/6a030751147a45c0863c3b5bde32c744/shares/d23b86ed-62d8-4fc2-b29a-e378fe1fa1fe/export_locations
