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

package cloudbackup

import (
	"strings"

	"github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/leveldbpool"
)

var (
	levelDBPool leveldbpool.LevelDB[uint]
	levelDBOnce once.Once
	LevelDBDir  = `data/cache/backup-db`
)

func LevelDB() leveldbpool.LevelDB[uint] {
	levelDBOnce.Do(initLevelDB)
	return levelDBPool
}

func initLevelDB() {
	levelDBPool = leveldbpool.New[uint](LevelDBDir)
}

func ParseDBValue(val []byte) (md5 string, startTs, endTs, fileModifyTs, fileSize int64) {
	parts := strings.Split(com.Bytes2str(val), `||`)
	md5 = parts[0]
	if len(parts) > 1 {
		com.SliceExtractCallback(parts[1:], func(v string) int64 {
			return param.AsInt64(v)
		}, &startTs, &endTs, &fileModifyTs, &fileSize)
	}
	return
}
