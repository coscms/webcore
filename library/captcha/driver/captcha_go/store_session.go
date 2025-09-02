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
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewStoreSession() captcha.Storer {
	return &storeSession{}
}

const (
	sessionKey = `captchaGoGen`
)

type storeSession struct {
}

func (a *storeSession) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	eCtx.Session().Set(sessionKey, com.String(val))
	return nil
}

func (a *storeSession) Get(ctx context.Context, key string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	sessVal, ok := eCtx.Session().Get(sessionKey).(string)
	if !ok {
		return captcha.ErrIllegalKey
	}
	*(value.(*[]byte)) = com.Str2bytes(sessVal)
	return nil
}

func (a *storeSession) Delete(ctx context.Context, key string) error {
	eCtx := ctx.(echo.Context)
	eCtx.Session().Delete(sessionKey)
	return nil
}
