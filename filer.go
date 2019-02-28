package main

import (
	"fmt"
	"time"

	"github.com/pepabo/go-netapp/netapp"
)

const (
	_url    = "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	version = "1.7"
)

type Filer struct {
	FilerBase
	Client *netapp.Client
	Share  *ProjectShareMap
}

type FilerBase struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewFiler(name, host, username, password string) *Filer {
	f := &Filer{
		FilerBase: FilerBase{
			Name:     name,
			Host:     host,
			Username: username,
			Password: password,
		},
	}
	f.Init()
	return f
}

func (f *Filer) Init() {
	url := fmt.Sprintf(_url, f.Host)
	opt := &netapp.ClientOptions{
		BasicAuthUser:     f.Username,
		BasicAuthPassword: f.Password,
		SSLVerify:         false,
		Timeout:           30 * time.Second,
	}
	f.Client = netapp.NewClient(url, version, opt)
}
