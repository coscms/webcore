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
	api := strings.Replace(cfg.ApiURL, `{ip}`, ip, -1)
	cli := restyclient.Classic()
	if len(cfg.ApiKey) > 0 {
		cli.SetAuthToken(cfg.ApiKey)
	}
	if cfg.ApiBasicAuth != nil && cfg.ApiBasicAuth.Username != `` && cfg.ApiBasicAuth.Password != `` {
		cli.SetBasicAuth(cfg.ApiBasicAuth.Username, cfg.ApiBasicAuth.Password)
	}
	if len(cfg.ApiHeaders) > 0 {
		cli.SetHeaders(cfg.ApiHeaders)
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
