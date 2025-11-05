package ip2region

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/admpub/log"
	syncOnce "github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/restyclient"
)

var (
	region4    *ip2region.Ip2Region
	region6    *ip2region.Ip2Region
	dict4File  string
	dict6File  string
	once4      syncOnce.Once
	once6      syncOnce.Once
	memoryMode bool
)

func init() {
	dict4File = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region_v4.xdb`
	dict6File = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region_v6.xdb`
	memoryMode = com.GetenvBool(`IP2REGION_MEMORY_MODE`, false)
	extend.Register(`ip2region`, func() interface{} {
		return &IP2RegionConfig{}
	})
}

func SetDict4File(f4 string) {
	dict4File = f4
	once4.Reset()
}

func SetDict6File(f6 string) {
	dict6File = f6
	once6.Reset()
}

func SetInstance4(new4Instance *ip2region.Ip2Region) {
	if region4 == nil {
		region4 = new4Instance
	} else {
		oldRegion4 := *region4
		*region4 = *new4Instance
		oldRegion4.Close()
	}
}

func SetInstance6(new6Instance *ip2region.Ip2Region) {
	if region6 == nil {
		region4 = new6Instance
	} else {
		oldRegion6 := *region6
		*region6 = *new6Instance
		oldRegion6.Close()
	}
}

func initialize4() (err error) {
	if region4 != nil {
		region4.Close()
	}
	region4, err = ip2region.New(dict4File, memoryMode)
	if err != nil {
		err = fmt.Errorf(`ip2region.New(%s) error: %w`, dict4File, err)
		log.Error(err)
	}
	return
}

func initialize6() (err error) {
	if region6 != nil {
		region6.Close()
	}
	region6, err = ip2region.New(dict6File, memoryMode)
	if err != nil {
		err = fmt.Errorf(`ip2region.New(%s) error: %w`, dict6File, err)
		log.Error(err)
	}
	return
}

func IsInitialized4() bool {
	return region4 != nil
}

func IsInitialized6() bool {
	return region6 != nil
}

// ErrIsInvalidIP 解析 IPv6 时会报这个错误
func ErrIsInvalidIP(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(err.Error(), `invalid ip address`)
}

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

func searchByLocalDict(cfg *IP2RegionConfig, ip string) (info ip2region.IpInfo, err error) {
	defer func() {
		if e := recover(); e != nil {
			panicErr := echo.NewPanicError(e, nil).Parse(15)
			log.Error(panicErr)
			err = fmt.Errorf(`%v`, e)
		}
	}()
	if net.ParseIP(ip).To4() != nil {
		once4.Do(func() {
			if cfg != nil && len(cfg.IPv4Dict) > 0 {
				SetDict4File(cfg.IPv4Dict)
			}
			err = initialize4()
		})
		if err != nil {
			return
		}
		info, err = region4.MemorySearch(ip)
		return
	}
	once6.Do(func() {
		if cfg != nil && len(cfg.IPv6Dict) > 0 {
			SetDict6File(cfg.IPv6Dict)
		}
		err = initialize6()
	})
	if err != nil {
		return
	}
	info, err = region6.MemorySearch(ip)
	return
}

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
	info, err := searchByLocalDict(cfg, ip)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(info)
}

type APIBasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type IP2RegionConfig struct {
	Mode         string            `json:"mode"`
	APIURL       string            `json:"apiUrl"`
	APIKey       string            `json:"apiKey,omitempty"`
	APIBasicAuth *APIBasicAuth     `json:"apiBasicAuth,omitempty"`
	APIHeaders   map[string]string `json:"apiHeaders,omitempty"`
	IPv4Dict     string            `json:"ipv4Dict,omitempty"`
	IPv6Dict     string            `json:"ipv6Dict,omitempty"`
}

func GetIP2RegionConfig() (cfg *IP2RegionConfig, ok bool) {
	cfg, ok = config.FromFile().Extend.Get(`ip2region`).(*IP2RegionConfig)
	return
}

func IPInfo(ip string) (info ip2region.IpInfo, err error) {
	if len(ip) == 0 {
		return
	}
	cfg, ok := GetIP2RegionConfig()
	if ok && cfg.Mode == `api` {
		return requestAPI(cfg, ip)
	}
	return searchByLocalDict(cfg, ip)
}

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

func IsZero(str string) bool {
	return len(str) == 0 || str == `0`
}

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
