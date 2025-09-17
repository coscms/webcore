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
	"net/url"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var (
	hpool = sync.Pool{
		New: func() interface{} {
			return echo.H{}
		},
	}

	urlValuesPool = sync.Pool{
		New: func() interface{} {
			return url.Values{}
		},
	}

	stringMapPool = sync.Pool{
		New: func() interface{} {
			return param.StringMap{}
		},
	}
)

func HPoolGet() echo.H {
	return hpool.Get().(echo.H)
}

func HPoolRelease(m echo.H) {
	for k := range m {
		delete(m, k)
	}

	hpool.Put(m)
}

func URLValuesPoolGet() url.Values {
	return urlValuesPool.Get().(url.Values)
}

func URLValuesPoolRelease(m url.Values) {
	for k := range m {
		delete(m, k)
	}

	urlValuesPool.Put(m)
}

func StringMapPoolGet() param.StringMap {
	return stringMapPool.Get().(param.StringMap)
}

func StringMapPoolRelease(m param.StringMap) {
	for k := range m {
		delete(m, k)
	}

	stringMapPool.Put(m)
}
