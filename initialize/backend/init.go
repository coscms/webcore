package backend

import (
	"net/url"
	"os"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo/subdomains"
)

func init() {
	config.AddConfigInitor(func(c *config.Config) {
		c.AddReloader(func(newConfig *config.Config) {
			c.Sys.ReloadRealIPConfig(&newConfig.Sys, route.IRegister().Echo().RealIPConfig())
		})
	})

	prefix := os.Getenv(`NGING_BACKEND_URL_PREFIX`)
	if len(prefix) > 0 {
		SetPrefix(prefix)
	}
	bootconfig.OnStart(0, start)

	config.OnKeySetSettings(`base.backendURL`, onChangeBackendURL)
}

func onChangeBackendURL(diff config.Diff) error {
	if !bootconfig.IsWeb() || !diff.IsDiff {
		return nil
	}
	oldURL, ok := diff.Old.(string)
	if !ok {
		return nil
	}
	newURL, ok := diff.New.(string)
	if !ok {
		return nil
	}
	if len(newURL) > 0 {
		u, err := url.Parse(newURL)
		if err != nil {
			return err
		}
		if len(u.Host) > 0 {
			subdomains.Default.Add(httpserver.KindBackend+`@`+u.Host, httpserver.Backend.Router.Echo())
		}
	}
	if len(oldURL) > 0 {
		u, err := url.Parse(oldURL)
		if err != nil {
			return err
		}
		if len(u.Host) > 0 {
			subdomains.Default.RemoveHost(u.Host)
		}
	}
	return nil
}
