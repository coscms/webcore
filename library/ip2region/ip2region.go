package ip2region

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/admpub/ip2region/v3/binding/golang/xdb"
	"github.com/admpub/log"
	syncOnce "github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/restyclient"
)

var (
	// dict4File specifies the file path for IPv4 database
	dict4File string

	// dict6File specifies the file path for IPv6 database
	dict6File string

	// once4 ensures thread-safe initialization of IPv4 database
	once4, once4Cancel = syncOnce.OnceValues[*ip2region.Ip2Region, error](initialize4)

	// once6 ensures thread-safe initialization of IPv6 database
	once6, once6Cancel = syncOnce.OnceValues[*ip2region.Ip2Region, error](initialize6)

	// memoryMode indicates whether to load the database into memory for faster lookups
	memoryMode bool
)

func init() {
	dict4File = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region.xdb`
	dict6File = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region_v6.xdb`
	memoryMode = com.GetenvBool(`IP2REGION_MEMORY_MODE`, false)
	extend.Register(`ip2region`, func() interface{} {
		return &IP2RegionConfig{}
	})
}

// SetDict4File sets the dictionary file path for the IPv4 dictionary and resets the initialization state.
func SetDict4File(f4 string) {
	dict4File = f4
	if once4Cancel != nil {
		once6Cancel()
	}
}

// SetDict6File sets the dictionary file path for the IPv6 dictionary and resets the initialization state.
func SetDict6File(f6 string) {
	dict6File = f6
	if once6Cancel != nil {
		once6Cancel()
	}
}

// initialize4 initializes the IPv4 IP2Region database.
// It closes any existing database connection first, then loads the new database from dict4File.
// Returns error if the database fails to load.
func initialize4() (region4 *ip2region.Ip2Region, err error) {
	var isMemoryMode bool
	dictFile := dict4File
	cfg, _ := GetIP2RegionConfig()
	if cfg != nil {
		if len(cfg.IPv4Dict) > 0 {
			dictFile = cfg.IPv4Dict
		}
		isMemoryMode = cfg.Mode == `local-memory`
	}
	if !isMemoryMode {
		isMemoryMode = memoryMode
	}
	region4, err = ip2region.New(dictFile, isMemoryMode)
	if err != nil {
		err = fmt.Errorf(`ip2region.New(%s) error: %w`, dictFile, err)
		log.Error(err)
	}
	return
}

// initialize6 initializes the IPv6 region database by loading the specified dictionary file.
// It closes any existing database connection before attempting to create a new one.
// Returns an error if the database initialization fails.
func initialize6() (region6 *ip2region.Ip2Region, err error) {
	var isMemoryMode bool
	dictFile := dict6File
	cfg, _ := GetIP2RegionConfig()
	if cfg != nil {
		if len(cfg.IPv6Dict) > 0 {
			dictFile = cfg.IPv6Dict
		}
		isMemoryMode = cfg.Mode == `local-memory`
	}
	if !isMemoryMode {
		isMemoryMode = memoryMode
	}
	region6, err = ip2region.New(dictFile, isMemoryMode)
	if err != nil {
		err = fmt.Errorf(`ip2region.New(%s) error: %w`, dictFile, err)
		log.Error(err)
	}
	return
}

// ErrIsInvalidIP 解析 IPv6 时会报这个错误
func ErrIsInvalidIP(err error) bool {
	if err == nil {
		return false
	}

	return xdb.IsInvalidIPAddress(err) || strings.HasPrefix(err.Error(), `invalid ip address`)
}

func ErrIsNotFoundXDB(err error) bool {
	if err == nil {
		return false
	}
	return os.IsNotExist(err)
}

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

// searchByLocalDict searches for IP location information using local dictionary files.
// It handles both IPv4 and IPv6 addresses, initializing the appropriate dictionary if needed.
// Returns IpInfo containing location details or an error if the search fails.
// Panics are recovered and logged, converting them to error returns.
func searchByLocalDict(ip string) (info ip2region.IpInfo, err error) {
	defer func() {
		if e := recover(); e != nil {
			panicErr := echo.NewPanicError(e, nil).Parse(15)
			log.Error(panicErr)
			err = fmt.Errorf(`%v`, e)
		}
	}()
	if net.ParseIP(ip).To4() != nil {
		var region4 *ip2region.Ip2Region
		region4, err = once4()
		if err != nil {
			return
		}
		info, err = region4.MemorySearch(ip)
		return
	}
	var region6 *ip2region.Ip2Region
	region6, err = once6()
	if err != nil {
		return
	}
	info, err = region6.MemorySearch(ip)
	return
}

// IP2RegionHandler handles IP address lookup requests using the IP2Region database.
// It validates the request IP parameter and checks API authentication if configured.
// Returns the region information in JSON format or an appropriate HTTP error:
// - 400 Bad Request if IP parameter is missing
// - 401 Unauthorized if authentication fails
// - 500 Internal Server Error if lookup fails
// - 200 OK with region data on success
func IP2RegionHandler(c echo.Context) error {
	ip := c.Param(`ip`)
	if len(ip) == 0 {
		return echo.ErrBadRequest
	}
	cfg, _ := GetIP2RegionConfig()
	if cfg != nil {
		if len(cfg.APIKey) > 0 {
			val := c.Header(echo.HeaderAuthorization)
			if strings.TrimPrefix(val, `Bearer `) != cfg.APIKey {
				return echo.ErrUnauthorized
			}
		} else if cfg.APIBasicAuth != nil && cfg.APIBasicAuth.Username != `` && cfg.APIBasicAuth.Password != `` {
			if username, password, ok := c.Request().BasicAuth(); !ok || username != cfg.APIBasicAuth.Username || password != cfg.APIBasicAuth.Password {
				c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Basic realm=Restricted")
				return echo.ErrUnauthorized
			}
		}
	}
	info, err := searchByLocalDict(ip)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(info)
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
	APIURL       string            `json:"apiUrl"`
	APIKey       string            `json:"apiKey,omitempty"`
	APIBasicAuth *APIBasicAuth     `json:"apiBasicAuth,omitempty"`
	APIHeaders   map[string]string `json:"apiHeaders,omitempty"`

	// 本地模式
	IPv4Dict string `json:"ipv4Dict,omitempty"`
	IPv6Dict string `json:"ipv6Dict,omitempty"`
}

// GetIP2RegionConfig retrieves the IP2Region configuration from the extended config file.
// Returns the configuration and a boolean indicating if the config was found and is of the correct type.
func GetIP2RegionConfig() (cfg *IP2RegionConfig, ok bool) {
	cfg, ok = config.FromFile().Extend.Get(`ip2region`).(*IP2RegionConfig)
	return
}

// IPInfo retrieves IP address information using either API mode or local dictionary search.
// It returns IpInfo struct containing location details and an error if any occurred.
// If the IP string is empty, returns zero values.
func IPInfo(ip string) (info ip2region.IpInfo, err error) {
	if len(ip) == 0 {
		return
	}
	cfg, ok := GetIP2RegionConfig()
	if ok && cfg.Mode == `api` {
		return requestAPI(cfg, ip)
	}
	return searchByLocalDict(ip)
}

// ClearZero removes "0" values from IpInfo fields by replacing them with empty strings.
// It checks and clears zero values for Country, Province, City, District, and ISP fields.
func ClearZero(info *ip2region.IpInfo) {
	if info.Country == `0` {
		info.Country = ``
	}
	if info.Province == `0` {
		info.Province = ``
	}
	if info.City == `0` {
		info.City = ``
	}
	if info.District == `0` {
		info.District = ``
	}
	if info.ISP == `0` {
		info.ISP = ``
	}
}

// IsZero checks if the string is empty or equals "0".
func IsZero(str string) bool {
	return len(str) == 0 || str == `0`
}

// Stringify converts IpInfo struct into a formatted string representation.
// If jsonify is true or not provided, returns a JSON-like string with fields like "国家", "省份", "城市".
// Otherwise returns a simple concatenated string of all non-zero fields.
// Only includes non-zero fields in the output.
func Stringify(info ip2region.IpInfo, jsonify ...bool) string {
	var (
		formats []string
		args    []interface{}
	)
	if !IsZero(info.Country) {
		formats = append(formats, `"国家":%q`)
		args = append(args, info.Country)
	}
	if !IsZero(info.Province) {
		formats = append(formats, `"省份":%q`)
		args = append(args, info.Province)
	}
	if !IsZero(info.City) {
		formats = append(formats, `"城市":%q`)
		args = append(args, info.City)
	}
	if !IsZero(info.District) {
		formats = append(formats, `"区县":%q`)
		args = append(args, info.District)
	}
	if !IsZero(info.ISP) {
		formats = append(formats, `"线路":%q`)
		args = append(args, info.ISP)
	}
	if len(jsonify) == 0 || jsonify[0] {
		return fmt.Sprintf(`{`+strings.Join(formats, `,`)+`}`, args...)
	}
	return fmt.Sprintf(strings.Repeat(`%s`, len(args)), args...)
}
