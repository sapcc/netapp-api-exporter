package netapp

import (
	"encoding/xml"
	"net/http"

	n "github.com/pepabo/go-netapp/netapp"
)

func (c *Client) ListCluster() (clusterInfo []n.ClusterIdentityInfo, resp *http.Response, err error) {
	opts := &n.ClusterIdentityOptions{}
	r, resp, err := c.ClusterIdentity.List(opts)
	clusterInfo = r.Results.ClusterIdentityInfo
	return
}

func (c *Client) CheckCluster() (statusCode int, err error) {
	body := c.ClusterIdentity
	body.Params.XMLName = xml.Name{Local: "cluster-identity-get"}
	resp, err := c.Do("POST", body)
	return resp.StatusCode, err
}

// func (c *Client) list() {
// 	body := c.ClusterIdentity
// 	body.Params.XMLName = xml.Name{Local: "cluster-identity-get"}
// 	r := n.ClusterIdentityResponse{}
// 	request, err := c.Client.NewRequest("POST", body)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	ctx, cncl := context.WithTimeout(context.Background(), c.ResponseTimeout)
// 	defer cncl()
// 	resp, err := checkResp(c.Client.Do(req.WithContext(ctx)))
// 	if err != nil {
// 		return nil, err
// 	}
// 	bs, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
// 	if c.options.Debug {
// 		log.Printf("[DEBUG] response xml \n%v\n", string(bs))
// 	}
// 	if v != nil {
// 		defer resp.Body.Close()
// 		err = xml.NewDecoder(resp.Body).Decode(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return resp, err
// }
