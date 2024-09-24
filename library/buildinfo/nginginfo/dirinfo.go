package nginginfo

import (
	"path/filepath"

	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/module"
	"github.com/webx-top/echo/middleware/render/driver"
)

func SetNgingDir(ngingDir string) {
	httpserver.Backend.AssetsDir = filepath.Join(ngingDir, backend.DefaultAssetsDir)
	httpserver.Backend.TemplateDir = filepath.Join(ngingDir, backend.DefaultTemplateDir)
}

func SetNgingPluginsDir(ngingPluginsDir string) {
	module.NgingPluginDir = ngingPluginsDir
}

func WatchTemplateDir(templateDirs ...string) {
	rendererDo := httpserver.Backend.RendererDo
	httpserver.Backend.RendererDo = func(renderer driver.Driver) {
		rendererDo(renderer)
		for _, templateDir := range templateDirs {
			renderer.Manager().AddWatchDir(templateDir)
		}
	}
}
