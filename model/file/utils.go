package file

import (
	"io"
	"mime"
	"strings"

	"github.com/coscms/go-imgparse/imgparse"
	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/echo"
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
