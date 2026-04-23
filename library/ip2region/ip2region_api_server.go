package ip2region

import (
	"net/http"
	"strings"

	"github.com/webx-top/echo"
)

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
	cfg, _ := GetIP2RegionConfig(c)
	if cfg != nil {
		if len(cfg.ApiKey) > 0 {
			val := c.Header(echo.HeaderAuthorization)
			if strings.TrimPrefix(val, `Bearer `) != cfg.ApiKey {
				return echo.ErrUnauthorized
			}
		} else if cfg.ApiBasicAuth != nil && cfg.ApiBasicAuth.Username != `` && cfg.ApiBasicAuth.Password != `` {
			if username, password, ok := c.Request().BasicAuth(); !ok || username != cfg.ApiBasicAuth.Username || password != cfg.ApiBasicAuth.Password {
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
