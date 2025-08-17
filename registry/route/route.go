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
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"

	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/route"
)

func IRegister() route.IRegister {
	return httpserver.Backend.Router
}

func Prefix() string {
	return IRegister().Prefix()
}

func SetPrefix(prefix string) {
	IRegister().SetPrefix(prefix)
}

func MakeHandler(handler interface{}, requests ...interface{}) echo.Handler {
	return IRegister().MakeHandler(handler, requests...)
}

func MetaHandler(handler interface{}, m echo.H, requests ...interface{}) echo.Handler {
	return IRegister().MetaHandler(m, handler, requests...)
}

func MetaHandlerWithRequest(handler interface{}, m echo.H, requests interface{}, methods ...string) echo.Handler {
	return IRegister().MetaHandlerWithRequest(m, handler, requests, methods...)
}

func HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return IRegister().HandlerWithRequest(handler, requests, methods...)
}

func Routes() []*echo.Route {
	return IRegister().Routes()
}

func Logger() logger.Logger {
	return IRegister().Logger()
}

func PreUse(middlewares ...interface{}) {
	IRegister().PreUse(middlewares...)
}

func PreToGroup(groupName string, middlewares ...interface{}) {
	IRegister().PreToGroup(groupName, middlewares...)
}

func Use(middlewares ...interface{}) {
	IRegister().Use(middlewares...)
}

// UseToGroup “@”符号代表后台网址前缀
func UseToGroup(groupName string, middlewares ...interface{}) {
	if groupName != `*` {
		groupName = `@` + groupName
	}
	IRegister().UseToGroup(groupName, middlewares...)
}

func AddGroupNamer(namers ...func(string) string) {
	IRegister().AddGroupNamer(namers...)
}

func Register(fn func(echo.RouteRegister)) {
	RegisterToGroup(``, fn)
}

func SetRootGroup(groupName string) {
	IRegister().SetRootGroup(groupName)
}

func Host(hostName string, middlewares ...interface{}) route.Hoster {
	return IRegister().Host(hostName, middlewares...)
}

func Clear() {
	IRegister().Clear()
}

func Apply() {
	echo.PanicIf(echo.Fire(`nging.route.apply.before`))
	IRegister().Apply()
	echo.PanicIf(echo.Fire(`nging.route.apply.after`))
}

func RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) route.MetaSetter {
	return IRegister().RegisterToGroup(`@`+groupName, fn, middlewares...)
}
