package ip2region

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/admpub/ip2region/v3/binding/golang/xdb"
	"github.com/admpub/log"
	syncOnce "github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
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
	cfg, _ := GetIP2RegionConfig(nil)
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
	cfg, _ := GetIP2RegionConfig(nil)
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

// IPInfo retrieves IP address information using either API mode or local dictionary search.
// It returns IpInfo struct containing location details and an error if any occurred.
// If the IP string is empty, returns zero values.
func IPInfo(c echo.Context, ip string) (info ip2region.IpInfo, err error) {
	if len(ip) == 0 {
		return
	}
	cfg, ok := GetIP2RegionConfig(c)
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
