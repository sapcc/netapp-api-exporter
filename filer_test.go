package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/pepabo/go-netapp/netapp"

	_ "github.com/motemen/go-loghttp/global"
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

func TestNetappClient(t *testing.T) {
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
	assert.NotNil(t, r.Results)
}

func TestNetappVolume(t *testing.T) {
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	az := os.Getenv("NETAPP_AZ")
	c := NewFiler("testFiler", host, username, password, az)

	opts := netapp.VolumeOptions{
		MaxRecords: 10,
	}

	volumePages := c.getNetappVolumePages(&opts, 1)
	vols := extracVolumes(volumePages)

	if assert.NotNil(t, vols) {
		fmt.Println("# of Vols: ", len(vols))

		for i, vol := range vols {
			fmt.Println("\nVolume: ", i)
			fmt.Println("Name\t\t", vol.VolumeIDAttributes.Name)
			fmt.Println("Type\t\t", vol.VolumeIDAttributes.Type)
			fmt.Println("Comment\t\t", vol.VolumeIDAttributes.Comment)
			fmt.Println("Node\t\t", vol.VolumeIDAttributes.Node)
			fmt.Println("Vserver\t\t", vol.VolumeIDAttributes.OwningVserverName)
			fmt.Println("AvailableSize\t", vol.VolumeSpaceAttributes.SizeAvailable)
			fmt.Println("Percentage\t", vol.VolumeSpaceAttributes.PercentageSizeUsed)
		}
	}
}

func TestManilaClient(t *testing.T) {
	var err error

	client, err := newManilaClient()
	if assert.Nil(t, err) {

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
			if assert.Nil(t, err) {
				if assert.NotNil(t, sel[0].ShareInstanceID) {
					fmt.Println("ShareInstanceID\t", sel[0].ShareInstanceID)
				}
			}
		}
	}

	if err != nil {
		assert.Equal(t, "", err.Error())
	}
}

func TestParseComment(t *testing.T) {
	var str, id, name, project string

	str = "share_id: 193b4209-2ef0-4752-a262-261b9fa27b25 in project: 631a3518e93d436fbdf57525babe8606"
	id, name, project = parseComment(str)
	assert.Equal(t, "193b4209-2ef0-4752-a262-261b9fa27b25", id)
	assert.Equal(t, "", name)
	assert.Equal(t, "631a3518e93d436fbdf57525babe8606", project)

	str = "share_id: 69fe1228-360c-4063-8f29-3a5bfb6d9772, share_name: c_blackbox_1553028005, project: d940aae3f8084f15a9b67de5b3b39720"
	id, name, project = parseComment(str)
	assert.Equal(t, "69fe1228-360c-4063-8f29-3a5bfb6d9772", id)
	assert.Equal(t, "c_blackbox_1553028005", name)
	assert.Equal(t, "d940aae3f8084f15a9b67de5b3b39720", project)
}
