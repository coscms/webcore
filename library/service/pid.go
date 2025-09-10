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

package service

import (
	stdLog "log"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

func createPidFile() string {
	pidFile := filepath.Join(com.SelfDir(), `data/pid`)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		stdLog.Println(err)
	}
	pidFile = filepath.Join(pidFile, `nging.pid`)
	return pidFile
}

func getPidFiles() []string {
	pidFile := []string{}
	pidFilePath := filepath.Join(com.SelfDir(), `data/pid`)
	err := filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == `daemon` { // 忽略进程值守创建的进程ID，避免被清理
				err = filepath.SkipDir
			}
			return err
		}
		if filepath.Ext(pidPath) == `.pid` {
			pidFile = append(pidFile, pidPath)
		}
		return nil
	})
	if err != nil {
		stdLog.Println(err)
	}
	return pidFile
}
