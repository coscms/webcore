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
	"context"
	"time"

	"github.com/webx-top/echo"
)

// New 新建后台执行信息
func New(c context.Context, opt echo.H) *Background {
	if c == nil {
		c = context.Background()
	}
	if opt == nil {
		opt = echo.H{}
	}
	ctx, cancel := context.WithCancel(c)
	return &Background{
		alone:   true,
		ctx:     ctx,
		cancel:  cancel,
		Options: opt,
		Started: time.Now(),
	}
}

// Background 后台执行信息
type Background struct {
	alone    bool
	op       string
	cacheKey string
	ctx      context.Context
	cancel   context.CancelFunc
	Options  echo.H
	Started  time.Time
}

// Context 暂存上下文信息
func (b *Background) Context() context.Context {
	return b.ctx
}

// Cancel 取消执行
func (b *Background) Cancel() {
	if b.alone {
		b.cancel()
		return
	}
	Cancel(b.op, b.cacheKey)
}
