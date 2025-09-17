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
	"errors"
	"html/template"
	"io/fs"
	"path"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/httpserver/httpserverutils"
)

type ICaptcha interface {
	Init(echo.H) error
	// keysValues: key1, value1, key2, value2
	Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML
	Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data
	MakeData(ctx echo.Context, hostAlias string, name string) echo.H
}

func RenderTemplate(ctx echo.Context, captchaType string, templatePath string, options param.Store) template.HTML {
	if len(templatePath) == 0 {
		templatePath = `default`
	}
	tmplPath, tmplFile := fixTemplatePath(captchaType, templatePath)
	b, err := ctx.Fetch(tmplPath, options)
	if err != nil {
		if templatePath != `default` && errors.Is(err, fs.ErrNotExist) {
			fileNotExist := true
			if httpserverutils.IsFrontendContext(ctx) && !strings.HasPrefix(templatePath, `#`) {
				b, err = ctx.Fetch(`#default#`+path.Join(`captcha`, captchaType, templatePath), options)
				fileNotExist = err != nil && errors.Is(err, fs.ErrNotExist)
			}
			if fileNotExist {
				tmplPath = strings.TrimSuffix(tmplPath, tmplFile)
				if !strings.HasSuffix(tmplPath, `/`) {
					tmplPath += `/`
				}
				b, err = ctx.Fetch(tmplPath+`default`, options)
			}
		}
		if err != nil {
			return template.HTML(err.Error())
		}
	}
	// return template.HTML(com.Bytes2str(b)) // bug: will be overwritten by a second render
	return template.HTML(string(b))
}
