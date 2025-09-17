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

package background

import (
	"sync"

	"github.com/webx-top/echo/param"
)

// Backgrounds 后台任务集合
var Backgrounds = sync.Map{}

// All 所有任务
func All() map[string]map[string]Background {
	r := map[string]map[string]Background{}
	Backgrounds.Range(func(key, val interface{}) bool {
		r[param.AsString(key)] = val.(*Group).Map()
		return true
	})
	return r
}

// ListBy 获取某个操作的所有任务
func ListBy(op string) *Group {
	old, exists := Backgrounds.Load(op)
	if !exists {
		return nil
	}
	exec := old.(*Group)
	return exec
}

// Cancel 取消执行
func Cancel(op string, cacheKeys ...string) {
	if len(cacheKeys) == 0 {
		return
	}
	exec := ListBy(op)
	if exec == nil {
		return
	}
	exec.Cancel(cacheKeys...)
	Backgrounds.Store(op, exec)
}
