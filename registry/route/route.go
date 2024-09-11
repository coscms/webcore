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

package route

import (
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/logger"

	"github.com/coscms/webcore/library/route"
)

var (
	routeRegister = route.NewRegister(defaults.Default).AddGroupNamer(groupNamer)
)

func init() {
	route.Default.Backend = routeRegister
}

func groupNamer(group string) string {
	if len(group) == 0 {
		return group
	}
	if group == `@` {
		return ``
	}
	if strings.HasPrefix(group, `@`) {
		return group[1:]
	}
	return group
}

func IRegister() route.IRegister {
	return routeRegister
}

func Prefix() string {
	return routeRegister.Prefix()
}

func SetPrefix(prefix string) {
	routeRegister.SetPrefix(prefix)
}

func MakeHandler(handler interface{}, requests ...interface{}) echo.Handler {
	return routeRegister.MakeHandler(handler, requests...)
}

func MetaHandler(handler interface{}, m echo.H, requests ...interface{}) echo.Handler {
	return routeRegister.MetaHandler(m, handler, requests...)
}

func MetaHandlerWithRequest(handler interface{}, m echo.H, requests interface{}, methods ...string) echo.Handler {
	return routeRegister.MetaHandlerWithRequest(m, handler, requests, methods...)
}

func HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return routeRegister.HandlerWithRequest(handler, requests, methods...)
}

func Routes() []*echo.Route {
	return routeRegister.Routes()
}

func Logger() logger.Logger {
	return routeRegister.Logger()
}

func Pre(middlewares ...interface{}) {
	routeRegister.Pre(middlewares...)
}

func PreToGroup(groupName string, middlewares ...interface{}) {
	routeRegister.PreToGroup(groupName, middlewares...)
}

func Use(middlewares ...interface{}) {
	routeRegister.Use(middlewares...)
}

// UseToGroup “@”符号代表后台网址前缀
func UseToGroup(groupName string, middlewares ...interface{}) {
	if groupName != `*` {
		groupName = `@` + groupName
	}
	routeRegister.UseToGroup(groupName, middlewares...)
}

func AddGroupNamer(namers ...func(string) string) {
	routeRegister.AddGroupNamer(namers...)
}

func Register(fn func(echo.RouteRegister)) {
	RegisterToGroup(``, fn)
}

func SetRootGroup(groupName string) {
	routeRegister.SetRootGroup(groupName)
}

func Host(hostName string, middlewares ...interface{}) route.Hoster {
	return routeRegister.Host(hostName, middlewares...)
}

func Clear() {
	routeRegister.Clear()
}

func Apply() {
	echo.PanicIf(echo.Fire(`nging.route.apply.before`))
	routeRegister.Apply()
	echo.PanicIf(echo.Fire(`nging.route.apply.after`))
}

func RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) route.MetaSetter {
	return routeRegister.RegisterToGroup(`@`+groupName, fn, middlewares...)
}
