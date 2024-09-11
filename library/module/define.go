package module

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
)

var NgingPluginDir = `vendor/github.com/nging-plugins` // `../../nging-plugins`
var WebCoreDir = `vendor/github.com/coscms/webcore`

func Register(modules ...IModule) {
	if len(modules) == 0 {
		return
	}
	schemaVer := echo.Float64(`SCHEMA_VER`)
	versionNumbers := []float64{schemaVer}
	for _, module := range modules {
		module.Apply()
		versionNumbers = append(versionNumbers, module.Version())
	}
	echo.Set(`SCHEMA_VER`, common.Float64Sum(versionNumbers...))
}
