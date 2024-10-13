package main

import (
	"github.com/coscms/webcore/library/module"
	"github.com/webx-top/echo"
)

var Module = &module.Module{
	DBSchemaVer: 0.1,
	Route:       registerRoute,
}

func registerRoute(r module.Router) {
	r.Backend().Register(func(r echo.RouteRegister) {
		r.Get(`/test`, Test)
	})
}

func Test(c echo.Context) error {
	return c.String(`test`)
}
