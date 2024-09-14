package module

import (
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/webx-top/echo/middleware"
)

type Misc struct {
	Template *ntemplate.PathAliases
	Assets   *middleware.StaticOptions
}
