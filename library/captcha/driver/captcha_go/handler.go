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
	"errors"
	"time"

	"github.com/coscms/captcha"
	"github.com/coscms/captcha/driver"
	"github.com/webx-top/echo"
)

func RegisterRoute(g echo.RouteRegister) {
	g.Get(`/:driver/:type`, CaptchaGoData)
	g.Post(`/:driver/:type`, CaptchaGoVerify)
	RegisterAPIRoute(g)
}

func CaptchaGoData(ctx echo.Context) error {
	resp := captcha.APIResponse{}
	c, err := driver.Singleton(ctx.Param(`driver`), ctx.Param(`type`), GetStorer)
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	data, err := c.MakeData(ctx)
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	resp.SetData(data)
	return ctx.JSON(resp)
}

func CaptchaGoVerify(ctx echo.Context) error {
	resp := captcha.APIResponse{}
	c, err := driver.Singleton(ctx.Param(`driver`), ctx.Param(`type`), GetStorer)
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	key := ctx.Form(`key`)
	err = c.Verify(ctx, key, ctx.Form(`response`))
	if err != nil {
		if errors.Is(err, captcha.ErrInvalidResponse) {
			c.Storer().Delete(ctx, key)
		}
		return ctx.JSON(resp.SetError(err.Error()))
	}
	c.Storer().Delete(ctx, key)
	captchaGoSetSuccessKey(ctx, key)
	return ctx.JSON(resp.SetSuccess())
}

const captchaGoSessionKey = `captchaGoKey`

func captchaGoSetSuccessKey(ctx echo.Context, key string) {
	secrets, ok := ctx.Session().Get(captchaGoSessionKey).(map[string]int64)
	if !ok {
		secrets = map[string]int64{}
	} else {
		// 始终只保存10条，避免session数据无限增长
		if len(secrets) >= 10 {
			var minTs int64
			var minKey string
			for k, v := range secrets {
				if minTs == 0 || minTs > v {
					minTs = v
					minKey = k
				}
			}
			delete(secrets, minKey)
		}
	}
	secrets[key] = time.Now().Unix()
	ctx.Session().Set(captchaGoSessionKey, secrets)
}

func captchaGoVerifySuccessKey(ctx echo.Context, key string, delAfterVerfiy bool) bool {
	secrets, ok := ctx.Session().Get(captchaGoSessionKey).(map[string]int64)
	if !ok {
		return ok
	}
	_, ok = secrets[key]
	if delAfterVerfiy && ok {
		delete(secrets, key)
		ctx.Session().Set(captchaGoSessionKey, secrets)
	}
	return ok
}
