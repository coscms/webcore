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

package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/nlog"
	"github.com/coscms/webcore/library/nlog/logcategory"
)

type Log struct {
	Debug        bool   `json:"debug"`
	Colorable    bool   `json:"colorable"`    // for console
	SaveFile     string `json:"saveFile"`     // for file
	FileMaxBytes int64  `json:"fileMaxBytes"` // for file
	Targets      string `json:"targets" form_delimiter:","`
}

func LogHasCategory() bool {
	return strings.Contains(FromFile().Settings().Log.LogFile(), `{category}`)
}

func (c *Log) Show(ctx echo.Context) error {
	category := ctx.Param(`category`, log.DefaultLog.Category)
	if strings.Contains(category, `..`) {
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`参数错误: %s`, category), 0).SetZone(`category`))
	}
	if category != log.DefaultLog.Category && !log.HasCategory(category) {
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`不存在日志分类: %s`, category), 0).SetZone(`category`))
	}
	logFile, err := logcategory.LogFilename(category, c.LogFile())
	if err != nil {
		return ctx.JSON(ctx.Data().SetError(err))
	}
	return nlog.LogShow(ctx, logFile)
}

func (c *Log) SetBy(r echo.H, defaults echo.H) *Log {
	if !r.Has(`log`) && defaults != nil {
		r.Set(`log`, defaults.GetStore(`log`))
	}
	loge := r.GetStore(`log`)
	c.Colorable = loge.Bool(`colorable`)
	c.SaveFile = loge.String(`saveFile`)
	switch t := loge.Get(`targets`).(type) {
	case []interface{}:
		for k, v := range t {
			if k > 0 {
				c.Targets += `,`
			}
			c.Targets = fmt.Sprint(v)
		}
	case []string:
		c.Targets = strings.Join(t, `,`)
	case string:
		c.Targets = t
	}
	c.FileMaxBytes = loge.Int64(`fileMaxBytes`)
	return c
}

func (c *Log) LogFile() string {
	if len(c.SaveFile) > 0 {
		return c.SaveFile
	}
	return filepath.Join(echo.Wd(), `data/logs/{category}_{date:20060102}_info.log`)
}

func (c *Log) Init() {
	//======================================================
	// 配置日志
	//======================================================
	if c.Debug {
		log.DefaultLog.MaxLevel = log.LevelDebug
		//log.DefaultLog.Formatter = log.ShortFileFormatter(0)
	} else {
		log.DefaultLog.MaxLevel = log.LevelInfo
	}
	targets := []log.Target{}
	for _, targetName := range strings.Split(c.Targets, `,`) {
		targetName = strings.TrimSpace(targetName)
		if len(targetName) == 0 {
			continue
		}
		switch targetName {
		case "file":
			//输出到文件
			logFileName := c.LogFile()
			if strings.Contains(logFileName, `{category}`) {
				for _, category := range logcategory.Categories.Slice() {
					fileTarget := log.NewFileTarget()
					fileTarget.Filter = &log.Filter{
						MaxLevel:   log.DefaultLog.MaxLevel,
						Categories: append([]string{category.K}, category.X...),
					}
					fileTarget.FileName = logFileName
					if c.FileMaxBytes > 0 {
						fileTarget.MaxBytes = c.FileMaxBytes
					}
					fileTarget.FileName = strings.Replace(logFileName, `{category}`, fileTarget.Filter.Categories[0], -1)
					targets = append(targets, fileTarget)
				}
			} else {
				fileTarget := log.NewFileTarget()
				fileTarget.Filter = &log.Filter{
					MaxLevel: log.DefaultLog.MaxLevel,
				}
				if c.FileMaxBytes > 0 {
					fileTarget.MaxBytes = c.FileMaxBytes
				}
				fileTarget.FileName = logFileName
				targets = append(targets, fileTarget)
			}

		case "console":
			fallthrough
		default:
			//输出到命令行
			consoleTarget := log.NewConsoleTarget()
			consoleTarget.ColorMode = c.Colorable
			targets = append(targets, consoleTarget)
		}
	}

	log.SetTarget(targets...)
	log.SetFatalAction(log.ActionExit)
}
