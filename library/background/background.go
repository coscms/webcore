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
	"sync/atomic"
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
		alone:    true,
		ctx:      ctx,
		cancel:   cancel,
		Options:  opt,
		Started:  time.Now(),
		Progress: &Progress{},
	}
}

type Progress struct {
	finish atomic.Int64
	total  atomic.Int64
}

func (b *Progress) Done(i int64) {
	b.finish.Add(i)
}

func (b *Progress) Reset() {
	b.finish.Store(0)
	b.total.Store(0)
}

func (b *Progress) SetFinish(i int64) {
	b.finish.Store(i)
}

func (b *Progress) SetTotal(i int64) {
	b.total.Store(i)
}

func (b *Progress) AddTotal(i int64) {
	b.total.Add(i)
}

func (b *Progress) Total() int64 {
	return b.total.Load()
}

func (b *Progress) Finish() int64 {
	return b.finish.Load()
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
	*Progress
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

func (b *Background) Clone() Background {
	prog := Progress{}
	prog.SetFinish(b.Progress.Finish())
	prog.SetTotal(b.Progress.Total())
	return Background{
		alone:    b.alone,
		op:       b.op,
		cacheKey: b.cacheKey,
		ctx:      b.ctx,
		cancel:   b.cancel,
		Options:  b.Options,
		Started:  b.Started,
		Progress: &prog,
	}
}
