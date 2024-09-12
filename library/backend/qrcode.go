package backend

import (
	"strings"

	"github.com/admpub/qrcode"
	"github.com/webx-top/echo"
)

func QrCode(ctx echo.Context) error {
	data := ctx.Form("data")
	size := ctx.Form("size")
	var (
		width  = 300
		height = 300
	)
	siz := strings.SplitN(size, `x`, 2)
	switch len(siz) {
	case 2:
		if i := ctx.Atop(siz[1]).Int(); i > 0 {
			height = i
		}
		fallthrough
	case 1:
		if i := ctx.Atop(siz[0]).Int(); i > 0 {
			width = i
		}
	}
	ctx.Response().Header().Set("Content-Type", "image/png")
	return qrcode.EncodeToWriter(data, width, height, ctx.Response())
}
