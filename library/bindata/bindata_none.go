//go:build !bindata
// +build !bindata

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

package bindata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/image"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/modal"
	uploadLibrary "github.com/coscms/webcore/library/upload"
)

// StaticOptions 后台static中间件选项
var StaticOptions = &middleware.StaticOptions{
	Root:     "",
	Path:     "",
	Fallback: []string{},
	MaxAge:   bootconfig.HTTPCacheMaxAge,
}

// Initialize 后台模板等素材初始化配置
func Initialize() {
	bootconfig.Bindata = false
	httpserver.Backend.StaticOptions = StaticOptions
	if !com.FileExists(bootconfig.FaviconPath) {
		log.Error(`not found favicon file: ` + bootconfig.FaviconPath)
	}
	bootconfig.FaviconHandler = func(c echo.Context) error {
		return c.CacheableFile(bootconfig.FaviconPath, bootconfig.HTTPCacheMaxAge)
	}
	image.WatermarkOpen = func(file string) (image.FileReader, error) {
		f, err := image.DefaultHTTPSystemOpen(file)
		if err != nil {
			if os.IsNotExist(err) {
				if strings.HasPrefix(file, uploadLibrary.UploadURLPath) || strings.HasPrefix(file, `/public/assets/`) {
					return os.Open(filepath.Join(echo.Wd(), file))
				}
			}
		}
		return f, err
	}
	modal.PathFixer = func(c echo.Context, file string) string {
		rpath := strings.TrimPrefix(file, httpserver.Backend.TemplateDir+`/`)
		rpath, ok := PathAliases.ParsePrefixOk(rpath)
		if ok {
			file = rpath
		}
		return file
	}
	httpserver.Backend.RendererDo = func(renderer driver.Driver) {
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			rpath, ok := PathAliases.ParsePrefixOk(tmpl)
			if ok {
				return rpath
			}
			tmpl = filepath.Join(renderer.TmplDir(), tmpl)
			return tmpl
		})
	}
}
