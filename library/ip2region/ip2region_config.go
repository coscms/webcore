package ip2region

import (
	"github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/webx-top/echo"
)

func init() {
	extend.Register(`ip2region`, func() interface{} {
		return &IP2RegionConfig{}
	})
}

// APIBasicAuth represents basic authentication credentials for API access.
// It contains Username and Password fields for authentication purposes.
type APIBasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// IP2RegionConfig defines the configuration structure for IP2Region service.
// It supports both API and local modes for IP address lookup.
// Fields:
//   - Mode: Operation mode ("api" or "local")
//   - APIURL: Endpoint URL for API mode
//   - APIKey: Authentication key for API mode (optional)
//   - APIBasicAuth: Basic authentication credentials for API mode (optional)
//   - APIHeaders: Custom headers for API requests (optional)
//   - IPv4Dict: Path to IPv4 database file for local mode (optional)
//   - IPv6Dict: Path to IPv6 database file for local mode (optional)
type IP2RegionConfig struct {
	Mode string `json:"mode"` // api / local / local-memory(or set env IP2REGION_MEMORY_MODE=true)

	// API 模式
	ApiURL       string            `json:"apiURL"` // 需要带{ip}占位符
	ApiKey       string            `json:"apiKey,omitempty"`
	ApiBasicAuth *APIBasicAuth     `json:"apiBasicAuth,omitempty"`
	ApiHeaders   map[string]string `json:"apiHeaders,omitempty"`

	// 本地模式
	IPv4Dict string `json:"ipv4Dict,omitempty"`
	IPv6Dict string `json:"ipv6Dict,omitempty"`
}

func (c *IP2RegionConfig) QueryRegion(ctx echo.Context, ip string) (info ip2region.IpInfo, err error) {
	if len(ip) == 0 {
		return
	}
	if c.Mode == `api` {
		return requestAPI(c, ip)
	}
	return searchByLocalDict(ip)
}

// GetIP2RegionConfig retrieves the IP2Region configuration from the extended config file.
// Returns the configuration and a boolean indicating if the config was found and is of the correct type.
func GetIP2RegionConfig(c echo.Context) (cfg *IP2RegionConfig, ok bool) {
	if c == nil {
		return getIP2RegionConfig()
	}

	cfg, ok = c.Internal().Get(`ip2region.config`).(*IP2RegionConfig)
	if ok {
		return
	}
	cfg, ok = getIP2RegionConfig()
	if !ok {
		return
	}
	c.Internal().Set(`ip2region.config`, cfg)
	return
}

func getIP2RegionConfig() (cfg *IP2RegionConfig, ok bool) {
	cfg, ok = config.FromDB(`thirdparty`).Get(`ip2region`).(*IP2RegionConfig)
	if ok && cfg.Mode != `` {
		return
	}
	cfg, ok = config.FromFile().Extend.Get(`ip2region`).(*IP2RegionConfig)
	return
}
