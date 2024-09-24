package dashboard

import "github.com/webx-top/echo"

var (
	_ echo.RenderContext = (*Tmplx)(nil)
	_ echo.RenderContext = (*Tmplxs)(nil)
)
