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
	"math"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// IDSharding 按照ID进行分片
func IDSharding(id uint64, shardingNum float64) uint64 {
	return uint64(math.Ceil(float64(id) / shardingNum))
}

// MD5Sharding 按照MD5进行分片
func MD5Sharding(str interface{}, length ...int) string {
	v := com.Md5(param.AsString(str))
	if len(length) == 0 || length[0] < 1 {
		return v[0:1]
	}
	if length[0] >= 32 {
		return v
	}
	return v[0:length[0]]
}

// MonthSharding 按照日期月进行分片
func MonthSharding(ctx echo.Context) string {
	return GetNowTime(ctx).Format("2006_01")
}

// YearSharding 按照日期年进行分片
func YearSharding(ctx echo.Context) string {
	return GetNowTime(ctx).Format("2006")
}

// GetNowTime 获取当前时间(同一个context中只获取一次)
func GetNowTime(ctx echo.Context) time.Time {
	t, y := ctx.Internal().Get(`time.now`).(time.Time)
	if !y || t.IsZero() {
		t = time.Now()
		ctx.Internal().Set(`time.now`, t)
	}
	return t
}

// DirShardingNum 文件夹分组基数
const DirShardingNum = float64(50000)

// DirSharding 文件夹分组(暂不使用)
func DirSharding(id uint64) uint64 {
	return IDSharding(id, DirShardingNum)
}
