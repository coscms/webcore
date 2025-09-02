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
	"html/template"

	"github.com/webx-top/echo"
	hdlCaptcha "github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var dflt ICaptcha = &defaultCaptcha{}

type defaultCaptcha struct {
}

func (c *defaultCaptcha) Init(_ echo.H) error {
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *defaultCaptcha) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("captchaId", GetHistoryOrNewCaptchaID(ctx, hdlCaptcha.DefaultOptions))
	if !options.Has("captchaName") {
		options.Set("captchaName", "code")
	}
	if len(templatePath) == 0 {
		return tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options)
	}
	options.Set("captchaImage", tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options))
	return RenderTemplate(ctx, TypeDefault, templatePath, options)
}

func (c *defaultCaptcha) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	return verifyAndSetDefaultCaptcha(ctx, hostAlias, name, captchaIdent...)
}

func (c *defaultCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	cid := GetHistoryOrNewCaptchaID(ctx, hdlCaptcha.DefaultOptions)
	return defaultCaptchaInfo(hostAlias, name, cid)
}
