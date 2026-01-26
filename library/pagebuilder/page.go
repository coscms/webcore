package pagebuilder

import (
	"github.com/webx-top/echo"
)

type Page struct {
	Tmpl          string
	Title         string
	Header        string
	Data          interface{}
	Body          echo.RenderContextWithData
	Footer        string
	Breadcrumb    []*Breadcrumb
	TopButtons    []*Button
	BottomButtons []*Button
}

func (p *Page) Render(ctx echo.Context) error {
	ctx.Set(`pageData`, p)
	if len(p.Tmpl) == 0 {
		switch p.Data.(type) {
		case *Table, Table:
			p.Tmpl = `common/page_table`
		case *Form, Form:
			p.Tmpl = `common/page_form`
		}
	}
	return ctx.Render(p.Tmpl, p.Data)
}
