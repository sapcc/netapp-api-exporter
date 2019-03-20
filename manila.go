package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	loghttp "github.com/motemen/go-loghttp"
)

type ManilaShare struct {
	ShareID       string
	ShareName     string
	ShareServerID string
	ProjectId     string
	InstanceID    string
}

func newManilaClient() (*gophercloud.ServiceClient, error) {
	region := os.Getenv("OS_REGION")
	identityEndpoint := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", region)

	client, err := openstack.NewClient(identityEndpoint)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{}
	config.InsecureSkipVerify = true

	var transport http.RoundTripper
	if os.Getenv("DEBUG") != "" {
		transport = &loghttp.Transport{
			Transport:  &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config},
			LogRequest: logHttpRequestWithHeader,
		}
	} else {
		transport = &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	}

	client.HTTPClient.Transport = transport

	opts := gophercloud.AuthOptions{
		DomainName: "ccadmin",
		TenantName: "cloud_admin",
		Username:   os.Getenv("OS_USERNAME"),
		Password:   os.Getenv("OS_PASSWORD"),
	}

	err = openstack.Authenticate(client, opts)
	if err != nil {
		logger.Printf("%+v", opts)
		return nil, err
	}

	eo := gophercloud.EndpointOpts{Region: region}

	manilaClient, err := openstack.NewSharedFileSystemV2(client, eo)
	if err != nil {
		return nil, err
	}

	manilaClient.Microversion = "2.46"
	return manilaClient, nil
}

func (f *Filer) GetManilaShare() (map[string]ManilaShare, error) {
	lo := shares.ListOpts{AllTenants: true}
	allpages, err := shares.ListDetail(f.OpenstackClient, lo).AllPages()
	if err != nil {
		return nil, err
	}

	sh, err := shares.ExtractShares(allpages)
	if err != nil {
		return nil, err
	}

	logger.Printf("%d Manila Shares fetched", len(sh))

	r := make(map[string]ManilaShare)
	for _, s := range sh {

		l, err := shares.GetExportLocations(f.OpenstackClient, s.ID).Extract()
		if err != nil {
			return nil, err
		}

		if len(l) > 0 {
			siid := l[0].ShareInstanceID
			siid = strings.Replace(siid, "-", "_", -1)
			r[siid] = ManilaShare{
				ShareID:       s.ID,
				ShareName:     s.Name,
				ShareServerID: s.ShareServerID,
				ProjectId:     s.ProjectID,
				InstanceID:    siid,
			}
		}
	}

	return r, nil
}

func logHttpRequestWithHeader(req *http.Request) {
	logger.Printf("--> %s %s %s", req.Method, req.URL, req.Header)
}
