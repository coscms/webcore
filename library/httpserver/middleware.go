package httpserver

import (
	"html/template"
	"net/url"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
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
			c.SetTransaction(factory.NewParam())
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
