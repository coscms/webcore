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

package captcha_go

import (
	"context"

	"github.com/coscms/captcha"
	"github.com/webx-top/echo"
)

func NewStoreCacheWrapper(cacher captcha.Storer) captcha.Storer {
	return &storeCacheWrapper{cacher: cacher}
}

type storeCacheWrapper struct {
	cacher captcha.Storer
}

func (a *storeCacheWrapper) Put(ctx context.Context, _ string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	a.cacher.Put(ctx, sessionKey+eCtx.Session().MustID(), val, timeout)
	return nil
}

func (a *storeCacheWrapper) Get(ctx context.Context, _ string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	return a.cacher.Get(ctx, sessionKey+eCtx.Session().MustID(), value)
}

func (a *storeCacheWrapper) Delete(ctx context.Context, _ string) error {
	eCtx := ctx.(echo.Context)
	return a.cacher.Delete(ctx, sessionKey+eCtx.Session().MustID())
}
