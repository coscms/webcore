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
)

var Multilingual = true

func NewRegister(e *echo.Echo, groupNamers ...func(string) string) IRegister {
	e.SetMultilingual(Multilingual)
	return &Register{
		echo:     e,
		handlers: Registers{},
		group:    NewGroup(groupNamers...),
		hosts:    make(map[string]*Host),
	}
}

type Register struct {
	echo           *echo.Echo
	prefix         string
	rootGroup      string
	handlers       Registers
	preMiddlewares []interface{}
	middlewares    []interface{}
	group          *Group
	hosts          map[string]*Host
	skipper        func() bool
}

func (r *Register) Clear() IRegister {
	r.echo = nil
	r.prefix = ``
	r.rootGroup = ``
	r.handlers = nil
	r.preMiddlewares = nil
	r.middlewares = nil
	r.group = nil
	r.hosts = nil
	r.skipper = nil
	return r
}

func (r *Register) Echo() *echo.Echo {
	return r.echo
}

func (r *Register) Routes() []*echo.Route {
	return r.echo.Routes()
}

func (r *Register) SetSkipper(f func() bool) IRegister {
	r.skipper = f
	return r
}

func (r *Register) Skipped() bool {
	if r == nil {
		return true
	}
	if r.skipper != nil {
		return r.skipper()
	}
	return false
}

func (r *Register) Logger() logger.Logger {
	return r.echo.Logger()
}

func (r *Register) Prefix() string {
	return r.prefix
}

func (r *Register) SetPrefix(prefix string) IRegister {
	r.prefix = prefix
	r.echo.SetPrefix(prefix)
	return r
}

func (r *Register) MetaHandler(m echo.H, handler interface{}, requests ...interface{}) echo.Handler {
	return r.echo.MetaHandler(m, handler, requests...)
}

func (r *Register) MakeHandler(handler interface{}, requests ...interface{}) echo.Handler {
	return r.echo.MakeHandler(handler, requests...)
}

func (r *Register) MetaHandlerWithRequest(m echo.H, handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return r.echo.MetaHandlerWithRequest(m, handler, requests, methods...)
}

func (r *Register) HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return r.echo.MetaHandlerWithRequest(nil, handler, requests, methods...)
}

func (r *Register) AddGroupNamer(namers ...func(string) string) IRegister {
	r.group.AddNamer(namers...)
	return r
}

func (r *Register) SetGroupNamer(namers ...func(string) string) IRegister {
	r.group.SetNamer(namers...)
	return r
}

func (r *Register) SetRootGroup(groupName string) IRegister {
	r.rootGroup = groupName
	return r
}

func (r *Register) RootGroup() string {
	return r.rootGroup
}

func (r *Register) Apply() IRegister {
	if r == nil || r.Skipped() {
		return r
	}
	e := r.echo
	e.Pre(r.preMiddlewares...)
	e.Use(r.middlewares...)
	r.handlers.Apply(e)
	r.group.Apply(e, r.rootGroup)
	for _, host := range r.hosts {
		hst := e.Host(host.Name, host.Middlewares...)
		if len(host.Alias) > 0 {
			hst.SetAlias(host.Alias)
		}
		host.Group.Apply(hst, r.rootGroup)
	}
	return r
}

func (r *Register) Pre(middlewares ...interface{}) IRegister {
	m := make([]interface{}, len(middlewares))
	copy(m, middlewares)
	r.preMiddlewares = append(m, r.preMiddlewares...)
	return r
}

func (r *Register) PreUse(middlewares ...interface{}) IRegister {
	r.preMiddlewares = append(r.preMiddlewares, middlewares...)
	return r
}

func (r *Register) Use(middlewares ...interface{}) IRegister {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Register) PreToGroup(groupName string, middlewares ...interface{}) IRegister {
	r.group.Pre(groupName, middlewares...)
	return r
}

func (r *Register) PreUseToGroup(groupName string, middlewares ...interface{}) IRegister {
	r.group.PreUse(groupName, middlewares...)
	return r
}

func (r *Register) UseToGroup(groupName string, middlewares ...interface{}) IRegister {
	r.group.Use(groupName, middlewares...)
	return r
}

func (r *Register) Register(fn func(echo.RouteRegister)) IRegister {
	r.handlers.Register(fn)
	return r
}

func (r *Register) RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) MetaSetter {
	r.group.Register(groupName, fn, middlewares...)
	return newMeta(groupName, r.group)
}

func (r *Register) Host(hostName string, middlewares ...interface{}) Hoster {
	host, ok := r.hosts[hostName]
	if !ok {
		host = &Host{
			Name:  hostName,
			Group: NewGroup(),
		}
		r.hosts[hostName] = host
	}
	host.Use(middlewares...)
	return host
}
