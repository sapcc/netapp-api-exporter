package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pepabo/go-netapp/netapp"
	"github.com/stretchr/testify/assert"
)

func init() {
	host := os.Getenv("NETAPP_HOST")
	url = fmt.Sprintf(_url, host)

	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
}

func TestB(t *testing.T) {
	c := netapp.NewClient(url, version, &netapp.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	})

	v := c.VolumeSpace
	p := netapp.VolumeSpaceOptions{
		MaxRecords: 10,
	}

	r, _, _ := v.List(&p)

	assert.True(t, r.Results.Passed())
	// assert.NotNil(t, r.Results)
}
