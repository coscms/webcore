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
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/service"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type LogCategories struct {
	WithCategory bool        `json:"withCategory"`
	Categories   echo.KVList `json:"categories"`
	OutputFile   string      `json:"outputFile"`
}

func (a *LogCategories) addCategory(k, v string) {
	logFilename, _ := LogFilename(k, a.OutputFile)
	if len(logFilename) > 0 {
		logFilename = filepath.Base(logFilename)
	}
	a.Categories.Add(k, v, echo.KVOptHKV(`logFilename`, logFilename))
}

func LogList(ctx echo.Context, hasCategory bool, logOutputFile string) LogCategories {
	logs := LogCategories{
		Categories: echo.KVList{},
		OutputFile: logOutputFile,
	}
	logs.WithCategory = hasCategory
	for _, v := range Categories.Slice() {
		if v.K == log.DefaultLog.Category {
			logs.addCategory(v.K, ctx.T(v.V))
			if !logs.WithCategory {
				break
			}
		}
		if len(v.V) == 0 {
			logs.addCategory(v.K, ctx.T(`%s日志`, com.Title(v.K)))
		} else {
			logs.addCategory(v.K, ctx.T(v.V))
		}
	}
	return logs
}

func LogFilename(category string, logOutputFile string) (string, error) {
	_, _, timeformat, filename, err := log.DateFormatFilename(logOutputFile)
	if err != nil {
		return ``, err
	}
	var logFile string
	if len(timeformat) > 0 {
		logFile = fmt.Sprintf(filename, time.Now().Format(timeformat))
	} else {
		logFile = filename
	}
	logFile = strings.ReplaceAll(logFile, `{category}`, category)
	if com.FileExists(logFile) {
		return logFile, nil
	}
	serviceAppLogFile := service.ServiceLogDir() + echo.FilePathSeparator + service.ServiceAppLogFile
	_, _, timeformat, filename, err = log.DateFormatFilename(serviceAppLogFile)
	if err == nil {
		if len(timeformat) > 0 {
			serviceAppLogFile = fmt.Sprintf(filename, time.Now().Format(timeformat))
		} else {
			serviceAppLogFile = filename
		}
		serviceAppLogFile = strings.ReplaceAll(serviceAppLogFile, `{category}`, category)
		if com.FileExists(serviceAppLogFile) {
			logFile = serviceAppLogFile
		}
	}
	return logFile, nil
}
