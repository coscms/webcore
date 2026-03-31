package file

import (
	"io"
	"mime"
	"strings"

	"github.com/coscms/go-imgparse/imgparse"
	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func ParseMime(m *dbschema.NgingFile, forceReset bool) {
	if forceReset || len(m.Mime) == 0 {
		m.Mime = mime.TypeByExtension(m.Ext)
		if len(m.Mime) == 0 {
			m.Mime = echo.MIMEOctetStream
		}
	}
}

func ParseImage(m *dbschema.NgingFile, reader io.Reader, forceReset bool) error {
	if m.Type != `image` {
		return nil
	}
	if forceReset || m.Width == 0 || m.Height == 0 {
		typ := strings.TrimPrefix(m.Ext, `.`)
		if typ == `jpg` {
			typ = `jpeg`
		}
		width, height, err := imgparse.Parse(reader, typ)
		if err != nil {
			return err
		}
		m.Width = uint(width)
		m.Height = uint(height)
		m.Dpi = 0
	}
	return nil
}

func MakeImageWidthAndHeightChecker(minWidth, maxWidth, minHeight, maxHeight uint) func(ctx echo.Context, nf *dbschema.NgingFile) error {
	return func(ctx echo.Context, nf *dbschema.NgingFile) error {
		if nf.Type != `image` {
			return nil
		}
		if minWidth > 0 && nf.Width < minWidth {
			return ctx.NewError(code.DataFormatIncorrect, `图片宽度不能小于%d像素`, minWidth)
		}
		if maxWidth > 0 && nf.Width > maxWidth {
			return ctx.NewError(code.DataFormatIncorrect, `图片宽度不能大于%d像素`, maxWidth)
		}
		if minHeight > 0 && nf.Height < minHeight {
			return ctx.NewError(code.DataFormatIncorrect, `图片高度不能小于%d像素`, minHeight)
		}
		if maxHeight > 0 && nf.Height > maxHeight {
			return ctx.NewError(code.DataFormatIncorrect, `图片高度不能大于%d像素`, maxHeight)
		}
		return nil
	}
}
