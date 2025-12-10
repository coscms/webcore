package config

import (
	"net/http"

	"github.com/coscms/webcore/library/formbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/language"
)

func init() {
	code.CodeDict.SetHTTPCodeToExists(code.NonPrivileged, http.StatusForbidden)
	formbuilder.LanguageConfigGetter = func(ctx echo.Context) language.Config {
		return FromFile().Language
	}
}
