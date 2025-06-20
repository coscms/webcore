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

package logcategory

import (
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func addLogCategory(logCategories *echo.KVList, k, v string) {
	logFilename, _ := config.FromFile().Settings().Log.LogFilename(k)
	if len(logFilename) > 0 {
		logFilename = filepath.Base(logFilename)
	}
	logCategories.Add(k, v, echo.KVOptHKV(`logFilename`, logFilename))
}

type LogCategories struct {
	WithCategory bool        `json:"withCategory"`
	Categories   echo.KVList `json:"categories"`
}

func LogList(ctx echo.Context) LogCategories {
	logs := LogCategories{}
	logCategories := &echo.KVList{}
	addLogCategory(logCategories, log.DefaultLog.Category, ctx.T(`Nging日志`))
	if strings.Contains(config.FromFile().Settings().Log.LogFile(), `{category}`) {
		logs.WithCategory = true
		categories := config.FromFile().Settings().Log.LogCategories()
		for _, k := range categories {
			k = strings.SplitN(k, `,`, 2)[0]
			v := k
			switch k {
			case `db`:
				v = ctx.T(`SQL日志`)
			case `echo`:
				v = ctx.T(`Web框架日志`)
			default:
				v = ctx.T(`%s日志`, com.Title(k))
			}
			addLogCategory(logCategories, k, v)
		}
	}
	logs.Categories = *logCategories
	return logs
}
