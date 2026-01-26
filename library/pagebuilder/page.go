package pagebuilder

import (
	"github.com/webx-top/echo"
)

type Page struct {
	Template      string
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
	if len(p.Template) == 0 {
		switch p.Data.(type) {
		case *Table, Table:
			p.Template = `common/page_table`
		case *Form, Form:
			p.Template = `common/page_form`
		}
	}
	return ctx.Render(p.Template, p.Data)
}
