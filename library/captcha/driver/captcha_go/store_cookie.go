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
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewStoreCookie() captcha.Storer {
	return &storeCookie{}
}

type storeCookie struct {
}

const (
	cookieKey = `CaptchaGo`
)

func (a *storeCookie) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	value := fmt.Sprintf(`%d|%s|%s`, time.Now().Unix(), com.Md5(eCtx.RealIP()+eCtx.Request().UserAgent()), val)
	log.Debugf(`set cookie value for %s: %s`, cookieKey, value)
	eCtx.Cookie().EncryptSet(cookieKey, value, captcha.MaxAge)
	return nil
}

func (a *storeCookie) Get(ctx context.Context, key string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	cookieVal := eCtx.Cookie().DecryptGet(cookieKey)
	if len(cookieVal) == 0 {
		log.Debugf(`failed to get cookie for %s`, cookieKey)
		return captcha.ErrIllegalKey
	}
	parts := strings.SplitN(cookieVal, `|`, 3)
	if len(parts) != 3 {
		log.Debugf(`failed to parse cookie for %s: %s`, cookieKey, cookieVal)
		return captcha.ErrIllegalKey
	}
	if time.Now().Unix()-param.AsInt64(parts[0]) > captcha.MaxAge {
		log.Debugf(`cookie has expired: %s`, cookieKey)
		return captcha.ErrIllegalKey
	}
	if parts[1] != com.Md5(eCtx.RealIP()+eCtx.Request().UserAgent()) {
		log.Debugf(`illegal cookie: [%s] %s`, eCtx.RealIP(), eCtx.Request())
		return captcha.ErrIllegalKey
	}
	*(value.(*[]byte)) = []byte(parts[2])
	return nil
}

func (a *storeCookie) Delete(ctx context.Context, key string) error {
	eCtx := ctx.(echo.Context)
	eCtx.Cookie().Set(cookieKey, ``, -1)
	return nil
}
