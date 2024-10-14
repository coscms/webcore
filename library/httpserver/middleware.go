package httpserver

import (
	"html/template"
	"net/url"
	"regexp"
	"strings"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/upload"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func MaxRequestBodySize(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().SetMaxSize(config.FromFile().GetMaxRequestBodySize())
		return h.Handle(c)
	}
}

func Transaction() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetTransaction(factory.NewParam().SetContext(c))
			return h.Handle(c)
		})
	}
}

var (
	EmptyURL = &url.URL{}
)

func ErrorPageFunc(c echo.Context) error {
	var siteURI *url.URL
	siteURL := config.Setting(`base`).String(`siteURL`)
	if len(siteURL) > 0 {
		siteURI, _ = url.Parse(siteURL)
	}
	c.Internal().Set(`siteURI`, siteURI)
	c.SetFunc(`SiteURI`, func() *url.URL {
		if siteURI == nil {
			return EmptyURL
		}
		return siteURI
	})
	c.SetFunc(`CaptchaForm`, func(tmpl string, args ...interface{}) template.HTML {
		return captchabiz.CaptchaForm(c, tmpl, args...)
	})
	c.SetFunc(`CaptchaFormWithType`, func(typ string, tmpl string, args ...interface{}) template.HTML {
		return captchabiz.CaptchaFormWithType(c, typ, tmpl, args...)
	})
	return nil
}

func TrimPathSuffix(ignorePrefixes ...string) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			upath := c.Request().URL().Path()
			for _, ignorePrefix := range ignorePrefixes {
				if strings.HasPrefix(upath, ignorePrefix) {
					return h.Handle(c)
				}
			}
			cleanedPath := strings.TrimSuffix(upath, c.DefaultExtension())
			c.Request().URL().SetPath(cleanedPath)
			return h.Handle(c)
		}
	}
}

func FixedUploadURLPrefix() echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			upath := c.Request().URL().Path()
			if strings.HasPrefix(upath, upload.UploadURLPath) {
				c.Request().URL().SetPath(c.Echo().Prefix() + upath)
			}
			return h.Handle(c)
		}
	}
}

func SearchEngineNoindex() echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			echo.SearchEngineNoindex(c)
			return h.Handle(c)
		}
	}
}

func HostChecker(key string) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			re, ok := echo.Get(key).(*regexp.Regexp)
			if !ok {
				return h.Handle(c)
			}
			// c.Host() 不含端口号
			if re.MatchString(c.Host()) {
				return h.Handle(c)
			}
			c.Response().NotFound()
			return nil
		}
	}
}
