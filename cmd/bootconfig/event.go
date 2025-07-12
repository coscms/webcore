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

package bootconfig

var onStartList []func()

// OnStart OnWebStart的别名函数。注册Web(管理面板)开始启动时执行的函数
func OnStart(index int, fn ...func()) {
	OnWebStart(index, fn...)
}

// OnWebStart 注册Web(管理面板)开始启动时执行的函数
func OnWebStart(index int, fn ...func()) {
	if len(fn) == 0 {
		return
	}
	if index < 0 {
		onStartList = append(onStartList, fn...)
		return
	}
	size := len(onStartList)
	if size > index {
		fn = append(fn, onStartList[index])
		onStartList[index] = fn[0]
		if len(fn) > 1 {
			OnStart(index+1, fn[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		onStartList = append(onStartList, nil)
	}
	onStartList = append(onStartList, fn...)
}

// WebStart Web(管理面板)开始启动
func WebStart() {
	for _, fn := range onStartList {
		if fn == nil {
			continue
		}
		fn()
	}
}
