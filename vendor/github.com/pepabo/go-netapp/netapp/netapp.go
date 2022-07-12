package netapp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	libraryVersion = "2.0.0"
	ServerURL      = `/servlets/netapp.servlets.admin.XMLrequest_filer`
	userAgent      = "go-netapp/" + libraryVersion
	XMLNs          = "http://www.netapp.com/filer/admin"
)

// A Client manages communication with the GitHub API.
type Client struct {
	client             *http.Client
	BaseURL            *url.URL
	UserAgent          string
	options            *ClientOptions
	ResponseTimeout    time.Duration
	Aggregate          *Aggregate
	AggregateSpace     *AggregateSpace
	AggregateSpares    *AggregateSpares
	Certificate        *Certificate
	Cf                 *Cf
	ClusterIdentity    *ClusterIdentity
	Diagnosis          *Diagnosis
	EnvironmentSensors *EnvironmentSensors
	Fcp                *Fcp
	Fcport             *Fcport
	Job                *Job
	Lun                *Lun
	Net                *Net
	Perf               *Perf
	Qtree              *Qtree
	QosPolicy          *QosPolicy
	Quota              *Quota
	QuotaReport        *QuotaReport
	QuotaStatus        *QuotaStatus
	Snapshot           *Snapshot
	Snapmirror         *Snapmirror
	StorageDisk        *StorageDisk
	System             *System
	Volume             *Volume
	VolumeSpace        *VolumeSpace
	VolumeOperations   *VolumeOperation
	LunOperations      *LunOperation
	VServer            *VServer
}

type ClientOptions struct {
	BasicAuthUser     string
	BasicAuthPassword string
	SSLVerify         bool
	CAFile            string
	CertFile          string
	KeyFile           string
	Debug             bool
	Timeout           time.Duration
}

func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		SSLVerify: true,
		Debug:     true,
		Timeout:   60 * time.Second,
	}
}

func NewClient(endpoint string, version string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if options.Timeout == 0 {
		options.Timeout = 60 * time.Second
	}

	tlsConfig, err := newTLSConfig(options)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: options.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	if !strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint + "/"
	}

	baseURL, _ := url.Parse(endpoint)

	c := &Client{
		client:          httpClient,
		BaseURL:         baseURL,
		UserAgent:       userAgent,
		options:         options,
		ResponseTimeout: options.Timeout,
	}

	b := Base{
		client:  c,
		XMLNs:   XMLNs,
		Version: version,
	}

	c.Aggregate = &Aggregate{
		Base: b,
	}

	c.AggregateSpace = &AggregateSpace{
		Base: b,
	}

	c.AggregateSpares = &AggregateSpares{
		Base: b,
	}

	c.Certificate = &Certificate{
		Base: b,
	}

	c.ClusterIdentity = &ClusterIdentity{
		Base: b,
	}
	c.Cf = &Cf{
		Base: b,
	}

	c.Diagnosis = &Diagnosis{
		Base: b,
	}

	c.EnvironmentSensors = &EnvironmentSensors{
		Base: b,
	}

	c.Fcp = &Fcp{
		Base: b,
	}

	c.Fcport = &Fcport{
		Base: b,
	}

	c.Job = &Job{
		Base: b,
	}

	c.Lun = &Lun{
		Base: b,
	}

	c.Net = &Net{
		Base: b,
	}

	c.Perf = &Perf{
		Base: b,
	}

	c.Qtree = &Qtree{
		Base: b,
	}

	c.QosPolicy = &QosPolicy{
		Base: b,
	}

	c.Quota = &Quota{
		Base: b,
	}

	c.QuotaReport = &QuotaReport{
		Base: b,
	}

	c.QuotaStatus = &QuotaStatus{
		Base: b,
	}

	c.Snapshot = &Snapshot{
		Base: b,
	}

	c.Snapmirror = &Snapmirror{
		Base: b,
	}

	c.StorageDisk = &StorageDisk{
		Base: b,
	}

	c.System = &System{
		Base: b,
	}

	c.Volume = &Volume{
		Base: b,
	}

	c.VolumeSpace = &VolumeSpace{
		Base: b,
	}

	c.VolumeOperations = &VolumeOperation{
		Base: b,
	}
	c.LunOperations = &LunOperation{
		Base: b,
	}

	c.VServer = &VServer{
		Base: b,
	}

	return c, nil
}

func (c *Client) NewRequest(method string, body interface{}) (*http.Request, error) {
	u, _ := c.BaseURL.Parse(ServerURL)

	buf, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, err
	}

	if c.options.Debug {
		log.Printf("[DEBUG] request xml: \n%v\n", string(buf))
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "text/xml")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	if c.options.BasicAuthUser != "" && c.options.BasicAuthPassword != "" {
		req.SetBasicAuth(c.options.BasicAuthUser, c.options.BasicAuthPassword)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	ctx, cncl := context.WithTimeout(context.Background(), c.ResponseTimeout)
	defer cncl()
	resp, err := checkResp(c.client.Do(req.WithContext(ctx)))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
	if c.options.Debug {
		log.Printf("[DEBUG] response xml \n%v\n", string(bs))
	}
	if v != nil {
		defer resp.Body.Close()
		err = xml.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}

// checkResp wraps an HTTP request from the default client and verifies that the
// request was successful. A non-200 request returns an error formatted to
// included any validation problems or otherwise.
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that.
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, newHTTPError(resp)
	}
}

func newHTTPError(resp *http.Response) error {
	return fmt.Errorf("Http Error status %d, Message: %s", resp.StatusCode, resp.Body)
}

// readCAFile reads the CA cert file from disk.
func readCAFile(f string) ([]byte, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("unable to load specified CA cert %s: %s", f, err)
	}
	return data, nil
}

func updateRootCA(cfg *tls.Config, b []byte) bool {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(b) {
		return false
	}
	cfg.RootCAs = caCertPool
	return true
}

// getClientCertificate reads the pair of client cert and key from disk and returns a tls.Certificate.
func getClientCertificate(options *ClientOptions) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(options.CertFile, options.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to use specified client cert (%s) & key (%s): %s", options.CertFile, options.KeyFile, err)
	}
	return &cert, nil
}

func newTLSConfig(options *ClientOptions) (*tls.Config, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: !options.SSLVerify}

	// If a CA cert is provided then let's read it in so we can validate the
	// scrape target's certificate properly.
	if len(options.CAFile) > 0 {
		b, err := readCAFile(options.CAFile)
		if err != nil {
			return nil, err
		}
		if !updateRootCA(tlsConfig, b) {
			return nil, fmt.Errorf("unable to use specified CA cert %s", options.CAFile)
		}
	}

	// If a client cert & key is provided then configure TLS config accordingly.
	if len(options.CertFile) > 0 && len(options.KeyFile) == 0 {
		return nil, fmt.Errorf("client cert file %q specified without client key file", options.CertFile)
	}

	if len(options.KeyFile) > 0 && len(options.CertFile) == 0 {
		return nil, fmt.Errorf("client key file %q specified without client cert file", options.KeyFile)
	}

	if len(options.CertFile) > 0 && len(options.KeyFile) > 0 {
		cert, err := getClientCertificate(options)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{*cert}
	}

	return tlsConfig, nil
}
