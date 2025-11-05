package sessionguard

import (
	"encoding/gob"

	"github.com/webx-top/echo"

	ip2regionparser "github.com/admpub/ip2region/v3/binding/golang/ip2region"
	"github.com/coscms/webcore/library/ip2region"
)

type Environment struct {
	UserAgent string                 `json:"ua"`
	Location  ip2regionparser.IpInfo `json:"loc"`
}

func init() {
	gob.Register(ip2regionparser.IpInfo{})
	gob.Register(&Environment{})
}

func (e *Environment) SetSession(ctx echo.Context, ownerType string) {
	ip2region.ClearZero(&e.Location)
	ctx.Session().Set(ownerType+`_env`, e)
}

func UnsetSession(ctx echo.Context, ownerType string) {
	ctx.Session().Delete(ownerType + `_env`)
}

func GetEnvFromSession(ctx echo.Context, ownerType string) *Environment {
	v, _ := ctx.Session().Get(ownerType + `_env`).(*Environment)
	return v
}
