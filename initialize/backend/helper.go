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
	"path/filepath"
	"sync"
	"time"

	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/render"

	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/captcha/driver/captcha_go"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/library/nretry"
	"github.com/coscms/webcore/middleware"
	"github.com/coscms/webcore/registry/route"
)

func Initialize() {
	route.Use(middleware.FuncMap(), middleware.BackendFuncMap(), render.Auto())
	route.Use(middleware.Middlewares...)
	addRouter()
	defaultConfigWatcher(true)
}

var onConfigChange = []func(file string) error{}

func OnConfigChange(fn func(file string) error) {
	onConfigChange = append(onConfigChange, fn)
}

func FireConfigChange(file string) error {
	err := nerrors.ErrIgnoreConfigChange
	for _, fn := range onConfigChange {
		if err := fn(file); err != nil {
			return err
		}
	}
	return err
}

var lockConfigChg = sync.Mutex{}

func defaultConfigWatcher(mustOk bool) {
	if config.FromCLI().Type != `manager` {
		return
	}
	conf := filepath.Base(config.FromCLI().Conf)
	config.WatchConfig(func(file string) error {
		lockConfigChg.Lock()
		defer lockConfigChg.Unlock()
		name := filepath.Base(file)
		switch name {
		case conf:
			time.Sleep(time.Second)
			err := nretry.OnErrorRetry(config.ParseConfig, 3, time.Second)
			if err != nil {
				if mustOk && config.IsInstalled() {
					config.MustOK(err)
				}
			}
			return err
		default:
			if !config.IsInstalled() {
				return nil
			}
			filePath := filepath.ToSlash(file)
			time.Sleep(time.Second)
			return FireConfigChange(filePath)
		}
	})
}

func addRouter() {
	captcha.New(``).Wrapper(route.IRegister().Echo()).SetMetaKV(httpserver.PermGuestKV())

	captchaGoG := route.IRegister().Echo().Group(`/captchago`, captchabiz.CheckEnable(captchaLib.TypeGo)).SetMetaKV(httpserver.PermGuestKV())
	captcha_go.RegisterRoute(captchaGoG)

	route.UseToGroup(`*`, middleware.AuthCheck) //应用中间件到所有子组
	route.Apply()
}
