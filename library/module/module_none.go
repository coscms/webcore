//go:build !bindata
// +build !bindata

package module

import (
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/httpserver"
)

func (m *Module) applyTemplateAndAssets() {
	m.setTemplate(bindata.PathAliases, httpserver.Frontend.TmplPathFixers.PathAliases)
	m.setAssets(bindata.StaticOptions, httpserver.Frontend.StaticOptions)
}

func SetBackendTemplate(key string, templatePath string) {
	SetTemplate(bindata.PathAliases, httpserver.Frontend.TmplPathFixers.PathAliases, key, templatePath)
}

func SetBackendAssets(assetsPath string) {
	SetAssets(bindata.StaticOptions, httpserver.Frontend.StaticOptions, assetsPath)
}
