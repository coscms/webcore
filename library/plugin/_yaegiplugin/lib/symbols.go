package lib

import "reflect"

//go:generate go run github.com/traefik/yaegi/cmd/yaegi@latest extract github.com/webx-top/echo
//go:generate go run github.com/traefik/yaegi/cmd/yaegi@latest extract github.com/webx-top/db
//go:generate go run github.com/traefik/yaegi/cmd/yaegi@latest extract github.com/coscms/webcore/library/module

var Symbols = map[string]map[string]reflect.Value{}
