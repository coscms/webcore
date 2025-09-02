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

package captcha

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var (
	//ErrCaptcha 验证码错误
	ErrCaptcha = echo.NewError(echo.T(`Captcha is incorrect`), code.CaptchaError)
	//ErrCaptchaIdMissing 缺少captchaId
	ErrCaptchaIdMissing = echo.NewError(echo.T(`Missing captchaId`), code.CaptchaIdMissing).SetZone(`captchaId`)
	//ErrCaptchaCodeRequired 验证码不能为空
	ErrCaptchaCodeRequired = echo.NewError(echo.T(`Captcha code is required`), code.CaptchaCodeRequired).SetZone(`code`)
)
