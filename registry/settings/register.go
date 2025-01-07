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

package settings

import (
	"fmt"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/errorslice"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var settings = []*SettingForm{
	{
		Short:    echo.T(`系统设置`),
		Label:    echo.T(`系统设置`),
		Group:    `base`,
		Tmpl:     []string{`manager/settings/base`},
		FootTmpl: []string{`manager/settings/base_footer`},
	},
	{
		Short: echo.T(`SMTP设置`),
		Label: echo.T(`SMTP服务器设置`),
		Group: `smtp`,
		Tmpl:  []string{`manager/settings/smtp`},
	},
	{
		Short: echo.T(`日志设置`),
		Label: echo.T(`日志设置`),
		Group: `log`,
		Tmpl:  []string{`manager/settings/log`},
	},
}

func Settings() []*SettingForm {
	return settings
}

func Register(sf ...*SettingForm) {
	registerSettings(sf...)
	for _, s := range sf {
		if s.items != nil {
			AddDefaultConfig(s.Group, s.items)
		}
		if s.dataEncoders != nil {
			s.dataEncoders.Register(s.Group)
		}
		if s.dataDecoders != nil {
			s.dataDecoders.Register(s.Group)
		}
	}
}

func registerSettings(sf ...*SettingForm) {
	for _, s := range sf {
		index, setting := Get(s.Group)
		if index == -1 {
			settings = append(settings, s)
		} else {
			settings[index] = setting.Merge(s)
		}
	}
}

func Get(group string) (int, *SettingForm) {
	for index, setting := range settings {
		if setting.Group == group {
			return index, setting
		}
	}
	return -1, nil
}

func RunHookPost(ctx echo.Context, groups ...string) error {
	n := len(groups)
	var i int
	errs := errorslice.New()
	for _, setting := range settings {
		if n < 1 || com.InSlice(setting.Group, groups) {
			err := setting.RunHookPost(ctx)
			if err != nil {
				err = fmt.Errorf("[config][group:%s] %s", setting.Group, err.(errorslice.Errors).ErrorTab())
				log.Error(err)
				errs.Add(err)
			}
			if n > 0 {
				i++
				if i >= n {
					break
				}
			}
		}
	}
	return errs.ToError()
}

func RunHookGet(ctx echo.Context, groups ...string) error {
	n := len(groups)
	var i int
	errs := errorslice.New()
	for _, setting := range settings {
		if n < 1 || com.InSlice(setting.Group, groups) {
			err := setting.RunHookGet(ctx)
			if err != nil {
				err = fmt.Errorf("[config][group:%s] %s", setting.Group, err.(errorslice.Errors).ErrorTab())
				log.Error(err)
				errs.Add(err)
			}
			if n > 0 {
				i++
				if i >= n {
					break
				}
			}
		}
	}
	return errs.ToError()
}
