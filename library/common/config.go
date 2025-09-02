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

package common

import (
	"strings"

	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config/subconfig/scookie"
)

const (
	ConfigName  = `ConfigFromFile`
	SettingName = `ConfigFromDB`
)

// APIKeyGetter API Key
type APIKeyGetter interface {
	APIKey() string
}

type ExtendConfigGetter interface {
	GetExtend() echo.H
}

type Cryptor interface {
	Encode(raw string, keys ...string) string
	Decode(encrypted string, keys ...string) string
	Encode256(raw string, keys ...string) string
	Decode256(encrypted string, keys ...string) string
}

type CookieConfigGetter interface {
	CookieConfig() scookie.Config
}

func CookieConfig() scookie.Config {
	return echo.Get(ConfigName).(CookieConfigGetter).CookieConfig()
}

func CryptorConfig() Cryptor {
	return echo.Get(ConfigName).(Cryptor)
}

func Setting(group ...string) echo.H {
	return echo.GetStoreByKeys(SettingName, group...)
}

func ExtendConfig() echo.H {
	return echo.Get(ConfigName).(ExtendConfigGetter).GetExtend()
}

func BackendURL(ctx echo.Context) string {
	backendURL := Setting(`base`).String(`backendURL`)
	if len(backendURL) == 0 {
		if ctx == nil {
			return backendURL
		}
		backendURL = ctx.Site()
	}
	backendURL = strings.TrimSuffix(backendURL, `/`)
	return backendURL
}

func SystemAPIKey() string {
	apiKey := Setting(`base`).String(`apiKey`)
	return apiKey
}
