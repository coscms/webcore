/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package backend

import (
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/pprof"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webcore/library/httpserver"
	ngingMW "github.com/coscms/webcore/middleware"
	"github.com/coscms/webcore/registry/route"
)

const (
	DefaultTemplateDir   = `./template/backend`      // 后台模板路径默认值
	DefaultAssetsDir     = `./public/assets/backend` // 后台素材路径默认值
	DefaultAssetsURLPath = `/public/assets/backend`  // 后台素材网址路径默认值
)

var (
	DefaultLocalHostNames = []string{
		`127.0.0.1`, `localhost`,
	}
)

func MakeSubdomains(domain string, appends []string) []string {
	domainList := strings.Split(domain, `,`)
	domain = domainList[0]
	if pos := strings.Index(domain, `://`); pos > 0 {
		pos += 3
		if pos < len(domain) {
			domain = domain[pos:]
		} else {
			domain = ``
		}
	}
	var myPort string
	domain, myPort = com.SplitHostPort(domain)
	if len(myPort) == 0 && len(domainList) > 1 {
		_, myPort = com.SplitHostPort(domainList[1])
	}
	port := strconv.Itoa(config.FromCLI().Port)
	newDomainList := []string{}
	if !com.InSlice(domain+`:`+port, domainList) {
		newDomainList = append(newDomainList, domain+`:`+port)
	}
	if myPort == port {
		myPort = ``
	}
	if len(myPort) > 0 {
		if !com.InSlice(domain+`:`+myPort, domainList) {
			newDomainList = append(newDomainList, domain+`:`+myPort)
		}
	}
	for _, hostName := range appends {
		if hostName == domain {
			continue
		}
		newDomainList = append(newDomainList, hostName+`:`+port)
		if len(myPort) > 0 {
			newDomainList = append(newDomainList, hostName+`:`+myPort)
		}
	}
	if len(newDomainList) > 0 {
		domainList = append(domainList, newDomainList...)
	}
	return param.StringSlice(domainList).Unique().String()
}

func SetPrefix(prefix string) {
	httpserver.Backend.SetPrefix(prefix)
}

func Prefix() string {
	return httpserver.Backend.Prefix()
}

func start() {
	e := httpserver.Backend.Router.Echo() // 不需要内部重启，所以直接操作*Echo
	config.FromFile().Sys.SetRealIPParams(e.RealIPConfig())
	e.SetRenderDataWrapper(echo.DefaultRenderDataWrapper)
	if len(Prefix()) > 0 {
		e.Pre(httpserver.FixedUploadURLPrefix())
	}

	// 子域名设置
	subdomains.Default.Default = httpserver.KindBackend
	subdomains.Default.Boot = httpserver.KindBackend
	domainName := subdomains.Default.Default
	backendDomain := config.FromCLI().BackendDomain
	if len(backendDomain) > 0 {
		domainName += `@` + strings.Join(MakeSubdomains(backendDomain, DefaultLocalHostNames), `,`)
	}
	subdomains.Default.Add(domainName, e)

	// 后台服务设置
	httpserver.Backend.GlobalFuncMap = backend.GlobalFuncMap()
	httpserver.Backend.Apply()
	httpserver.Backend.Renderer().MonitorEvent(func(file string) {
		if strings.HasSuffix(file, `.form.json`) {
			if formbuilder.DelCachedConfig(file) {
				log.Debug(`delete: cache form config: `, file)
			}
		}
	})
	//RendererDo(renderOptions.Renderer())
	echo.OnCallback(`nging.renderer.cache.clear`, func(_ events.Event) error {
		log.Debug(`clear: Backend Template Object Cache`)
		httpserver.Backend.Renderer().ClearCache()
		formbuilder.ClearCache()
		return nil
	})
	e.Get(`/favicon.ico`, bootconfig.FaviconHandler).SetMetaKV(route.PermGuestKV())
	httpserver.Backend.I18n().Handler(e, `App.i18n`)
	debugG := e.Group(`/debug`, ngingMW.DebugPprof).SetMetaKV(route.PermGuestKV())
	pprof.RegisterRoute(debugG)
	Initialize()
}
