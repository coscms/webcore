package backend

import (
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var onInstalled []func(ctx echo.Context) error

func OnInstalled(cb func(ctx echo.Context) error) {
	if cb == nil {
		return
	}
	onInstalled = append(onInstalled, cb)
}

func FireInstalled(ctx echo.Context) error {
	var err error
	for _, cb := range onInstalled {
		log.Info(color.GreenString(`[installer]`), `Execute Hook: `, com.FuncName(cb))
		if err = cb(ctx); err != nil {
			return ctx.NewError(code.Failure, err.Error())
		}
	}
	return err
}
