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

package captcha

import "sort"

const (
	TypeDefault = `default`
	TypeAPI     = `api`
	TypeGo      = `go`
)

var drivers = map[string]func() ICaptcha{
	TypeDefault: func() ICaptcha { return dflt },
}

func Register(name string, ic func() ICaptcha) {
	drivers[name] = ic
}

func Get(name string) func() ICaptcha {
	return drivers[name]
}

func GetOk(name string) (func() ICaptcha, bool) {
	ic, ok := drivers[name]
	return ic, ok
}

func GetAllNames() []string {
	names := make([]string, 0, len(drivers))
	for name := range drivers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func Has(name string) bool {
	_, ok := drivers[name]
	return ok
}
