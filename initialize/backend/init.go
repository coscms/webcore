package backend

import (
	"os"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/registry/route"
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
}
