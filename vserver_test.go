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

func TestA(t *testing.T) {
	c := netapp.NewClient(url, version, &netapp.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	})

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
