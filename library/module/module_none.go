//go:build !bindata
// +build !bindata

package module

import (
	"github.com/coscms/webcore/library/bindata"
)

var FrontendMisc = &Misc{}

func (m *Module) applyTemplateAndAssets() {
	m.setTemplate(bindata.PathAliases, FrontendMisc.Template)
	m.setAssets(bindata.StaticOptions, FrontendMisc.Assets)
}

func SetBackendTemplate(key string, templatePath string) {
	SetTemplate(bindata.PathAliases, FrontendMisc.Template, key, templatePath)
}

func SetBackendAssets(assetsPath string) {
	SetAssets(bindata.StaticOptions, FrontendMisc.Assets, assetsPath)
}
