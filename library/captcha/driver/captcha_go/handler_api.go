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
	"net/http"
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func RegisterAPIRoute(e echo.RouteRegister) {
	e.Route(`POST,GET,DELETE`, `/store`, storeAPIHandler, checkEnableAPIService)
}

type storeAPIRequest struct {
	Key     string `json:"key" validate:"required"`
	Val     string `json:"val,omitempty"`
	Timeout int64  `json:"timeout,omitempty"`
}

func storeAPIHandler(ctx echo.Context) error {
	data := ctx.Data()
	req := &storeAPIRequest{Key: ctx.Query(`key`)}
	err := ctx.MustBindAndValidate(req)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	err = req.checkToken(ctx)
	if err != nil {
		return ctx.JSON(data.SetError(err).SetCode(code.Unauthenticated.Int()))
	}
	storer, err := GetStorer()
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	switch ctx.Method() {
	case http.MethodGet:
		var val []byte
		err = storer.Get(ctx, req.Key, &val)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		data.SetData(string(val))
	case http.MethodDelete:
		err = storer.Delete(ctx, req.Key)
	case http.MethodPost:
		err = storer.Put(ctx, req.Key, req.Val, req.Timeout)
	}
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	return ctx.JSON(data.SetInfo(`success`, code.Success.Int()))
}

func (data *storeAPIRequest) checkToken(ctx echo.Context) error {
	apiKey := config.Setting(`base`).String(`apiKey`)
	if len(apiKey) == 0 {
		return echo.ErrNotAcceptable
	}
	token := ctx.Header(`Authorization`)
	parts := strings.SplitN(token, ` `, 2)
	if len(parts) == 2 {
		token = parts[1]
	} else {
		token = parts[0]
	}
	b, _ := com.JSONEncode(data)
	expectedToken := com.Token(apiKey, b)
	if expectedToken != token {
		return echo.ErrBadRequest
	}
	return nil
}

func checkEnableAPIService(h echo.Handler) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cfg := config.FromFile().Extend.GetStore(`captchaGo`)
		if !cfg.Bool(`apiService`) {
			return echo.ErrNotAcceptable
		}
		if cfg.String(`store`) == `api` { // 如果自己的接口请求自己的服务则会导致死循环
			return echo.ErrFailedDependency
		}
		return h.Handle(ctx)
	}
}
