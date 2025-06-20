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

package middleware

import (
	"html/template"
	"strings"

	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/modal"
	navigateLib "github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/nlog/logcategory"
	"github.com/coscms/webcore/library/role"
	"github.com/coscms/webcore/library/role/roleutils"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/coscms/webcore/registry/navigate"
	"github.com/coscms/webcore/registry/settings"
)

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Modal`, func(data interface{}) template.HTML {
				return modal.Render(c, data)
			})
			httpserver.ErrorPageFunc(c)
			return h.Handle(c)
		})
	}
}

func BackendFuncMap() echo.MiddlewareFunc {
	var emptyLogCategories logcategory.LogCategories
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {

			//用户相关函数
			user := backend.User(c)
			if user != nil {
				c.Set(`user`, user)
				c.SetFunc(`Username`, func() string { return user.Username })
				c.SetFunc(`IsFounder`, func() bool { return role.IsFounder(user) })
				c.Set(`roleList`, roleutils.UserRoles(c))
			}
			c.SetFunc(`ProjectIdent`, func() string {
				return GetProjectIdent(c)
			})
			c.SetFunc(`TopButtons`, func() dashboard.Buttons {
				buttons := httpserver.Backend.Dashboard.TopButtons.All(c)
				buttons.Ready(c)
				return buttons
			})
			c.SetFunc(`GlobalHeads`, func() dashboard.GlobalHeads {
				heads := httpserver.Backend.Dashboard.GlobalHeads.All(c)
				heads.Ready(c)
				return heads
			})
			c.SetFunc(`GlobalFooters`, func() dashboard.GlobalFooters {
				footers := httpserver.Backend.Dashboard.GlobalFooters.All(c)
				footers.Ready(c)
				return footers
			})
			c.SetFunc(`DashboardConfig`, func(extendOrType string) interface{} {
				// extendOrType:
				// 1. <extend>.<type[#buttonGroup]>
				// 2. <type[#buttonGroup]>
				parts := strings.SplitN(extendOrType, `.`, 2)
				var extend string
				var dtype string
				if len(parts) == 2 {
					extend = parts[0]
					dtype = parts[1]
				} else {
					dtype = parts[0]
				}
				var d *dashboard.Dashboard
				if len(extend) > 0 {
					d = httpserver.Backend.Dashboard.GetExtend(extend)
				} else {
					d = httpserver.Backend.Dashboard
				}
				if d == nil {
					return nil
				}
				return d.Get(c, dtype)
			})
			c.SetFunc(`SettingFormRender`, func(s *settings.SettingForm) interface{} {
				return s.Render(c)
			})
			c.SetFunc(`PermissionCheckByType`, func(permission role.ICheckByType, typ string, permPath string) interface{} {
				return permission.CheckByType(c, typ, permPath)
			})
			c.SetFunc(`Navigate`, func(side string) navigateLib.List {
				return GetBackendNavigate(c, side)
			})
			c.SetFunc(`HasNavigate`, func(navList *navigateLib.List) bool {
				if user != nil && role.IsFounder(user) {
					return true
				}
				permission := UserPermission(c)
				return permission.HasNavigate(c, navList)
			})
			c.SetFunc(`EnvKey`, func() string {
				return sessionguard.EnvKey(c, sessionguard.GetConfig().SessionGuardConfig)
			})
			c.SetFunc(`LogCategories`, func() logcategory.LogCategories {
				if user != nil && role.IsFounder(user) {
					return logcategory.LogList(c)
				}
				permission := UserPermission(c)
				if !permission.Check(c, `manager/log/:category`) {
					return emptyLogCategories
				}
				return logcategory.LogList(c)
			})
			return h.Handle(c)
		})
	}
}

func UserPermission(c echo.Context) *role.RolePermission {
	permission, ok := c.Internal().Get(`userPermission`).(*role.RolePermission)
	if !ok || permission == nil {
		permission = role.NewRolePermission().Init(roleutils.UserRoles(c))
		c.Internal().Set(`userPermission`, permission)
	}
	return permission
}

func GetProjectIdent(c echo.Context) string {
	projectIdent := c.Internal().String(`projectIdent`)
	if len(projectIdent) == 0 {
		projectIdent = navigate.ProjectIdent(c.Path())
		if len(projectIdent) == 0 {
			if proj := navigate.ProjectFirst(true); proj != nil {
				projectIdent = proj.Ident
			}
		}
		c.Internal().Set(`projectIdent`, projectIdent)
	}
	return projectIdent
}

func GetBackendNavigate(c echo.Context, side string) navigateLib.List {
	switch side {
	case `top`:
		navList, ok := c.Internal().Get(`navigate.top`).(navigateLib.List)
		if ok {
			return navList
		}
		user := backend.User(c)
		if user != nil && role.IsFounder(user) {
			if navigate.TopNavigate == nil {
				return navigate.EmptyList
			}
			return *navigate.TopNavigate
		}
		permission := UserPermission(c)
		navList = permission.FilterNavigate(c, navigate.TopNavigate)
		c.Internal().Set(`navigate.top`, navList)
		return navList
	case `left`:
		fallthrough
	default:
		navList, ok := c.Internal().Get(`navigate.left`).(navigateLib.List)
		if ok {
			return navList
		}
		user := backend.User(c)
		var leftNav *navigateLib.List
		ident := GetProjectIdent(c)
		if len(ident) > 0 {
			if proj := navigate.ProjectGet(ident); proj != nil {
				leftNav = proj.NavList
			}
		}
		if user != nil && role.IsFounder(user) {
			if leftNav == nil {
				return navigate.EmptyList
			}
			return *leftNav
		}
		permission := UserPermission(c)
		navList = permission.FilterNavigate(c, leftNav)
		c.Internal().Set(`navigate.left`, navList)
		return navList
	}
}
