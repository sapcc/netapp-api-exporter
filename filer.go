package main

import (
	"fmt"
	"time"

	"github.com/pepabo/go-netapp/netapp"
)

type Filer struct {
	FilerBase
	NetappClient *netapp.Client
}

type FilerBase struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	AvailabilityZone string `yaml:"availability_zone"`
}

func NewFiler(name, host, username, password, az string) *Filer {
	f := &Filer{
		FilerBase: FilerBase{
			Name:             name,
			Host:             host,
			Username:         username,
			Password:         password,
			AvailabilityZone: az,
		},
		NetappClient: newNetappClient(host, username, password),
	}
	f.Init()
	return f
}

func newNetappClient(host, username, password string) *netapp.Client {
	_url := "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	url := fmt.Sprintf(_url, host)

	version := "1.7"

	opts := &netapp.ClientOptions{
		BasicAuthUser:     username,
		BasicAuthPassword: password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}

	return netapp.NewClient(url, version, opts)
}
