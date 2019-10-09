package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pepabo/go-netapp/netapp"
)

type FilerManager struct {
	Filer
	NetappClient  *netapp.Client
	Volumes       []*NetappVolume
	Aggregates    []*NetappAggregate
	volMaxAge     time.Duration
	aggMaxAge     time.Duration
	lastVolScrape time.Time
	lastAggScrape time.Time
	mtxVol        sync.Mutex // protects lastVolScrape and Volumes
	mtxAgg        sync.Mutex // protects lastAggScrape and Aggregates
}

type Filer struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	AvailabilityZone string `yaml:"availability_zone"`
}

func NewFilerManager(f Filer) *FilerManager {
	return &FilerManager{
		Filer:        f,
		volMaxAge:    5 * time.Minute,
		aggMaxAge:    5 * time.Minute,
		NetappClient: newNetappClient(f.Host, f.Username, f.Password),
	}
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
