package main

import (
	"fmt"
	"time"

	"github.com/pepabo/go-netapp/netapp"
)

type NetappFilerClient struct {
	NetappFiler
	NetappClient *netapp.Client
}

type NetappFiler struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	AvailabilityZone string `yaml:"availability_zone"`
}

func NewNetappClient(f NetappFiler) NetappFilerClient {
	return NetappFilerClient{
		NetappFiler:  f,
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

func (f *NetappFilerClient) QueryAggregates(opts *netapp.AggrOptions) (res []netapp.AggrInfo, err error) {
	pageHandler := func(r netapp.AggrListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AggrAttributes...)
		return true
	}
	f.NetappClient.Aggregate.ListPages(opts, pageHandler)
	return
}

func (f *NetappFilerClient) QueryVolumes(opts *netapp.VolumeOptions) (res []netapp.VolumeInfo, err error) {
	pageHandler := func(r netapp.VolumeListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AttributesList...)
		return true
	}
	f.NetappClient.Volume.ListPages(opts, pageHandler)
	return
}
