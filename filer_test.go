package main

import (
	"os"
	"testing"

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
