package ip2region

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/webx-top/com"
	"github.com/webx-top/restyclient"
)

// requestAPI makes a request to the IP geolocation API with the given IP address.
// It handles API authentication (token or basic auth), custom headers, and response parsing.
// Returns IpInfo on success or error if the request fails or API returns non-success status.
func requestAPI(cfg *IP2RegionConfig, ip string) (ip2region.IpInfo, error) {
	api := strings.Replace(cfg.APIURL, `{ip}`, ip, -1)
	cli := restyclient.Classic()
	if len(cfg.APIKey) > 0 {
		cli.SetAuthToken(cfg.APIKey)
	}
	if cfg.APIBasicAuth != nil && cfg.APIBasicAuth.Username != `` && cfg.APIBasicAuth.Password != `` {
		cli.SetBasicAuth(cfg.APIBasicAuth.Username, cfg.APIBasicAuth.Password)
	}
	if len(cfg.APIHeaders) > 0 {
		cli.SetHeaders(cfg.APIHeaders)
	}
	resp, err := cli.Get(api)
	if err != nil {
		return ip2region.IpInfo{}, err
	}
	if !resp.IsSuccess() {
		return ip2region.IpInfo{}, fmt.Errorf(`IP geolocation API error: %s`, com.StripTags(resp.String()))
	}
	var info ip2region.IpInfo
	body := resp.Body()
	if bytes.HasPrefix(body, []byte(`{`)) && bytes.HasSuffix(body, []byte(`}`)) {
		err = json.Unmarshal(body, &info)
	} else {
		info.Parse(string(body))
	}
	return info, nil
}
