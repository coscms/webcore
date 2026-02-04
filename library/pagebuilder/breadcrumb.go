package pagebuilder

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewBreadcrumb(ctx echo.Context, title string, path string) *Breadcrumb {
	return &Breadcrumb{
		ctx:   ctx,
		Title: title,
		Path:  path,
	}
}

type Breadcrumb struct {
	ctx   echo.Context
	Icon  string
	Text  string
	Title string
	Path  string
}

func (b *Breadcrumb) URI() string {
	return com.WithURLParams(b.Path, `next`, b.ctx.RequestURI())
}

func (b *Breadcrumb) SetPath(path string) *Breadcrumb {
	b.Path = path
	return b
}

func (b *Breadcrumb) SetText(text string) *Breadcrumb {
	b.Text = text
	return b
}

func (b *Breadcrumb) SetTitle(title string) *Breadcrumb {
	b.Title = title
	return b
}

func (b *Breadcrumb) SetIcon(icon string) *Breadcrumb {
	b.Icon = icon
	return b
}
func (b *Breadcrumb) SetContext(ctx echo.Context) *Breadcrumb {
	b.ctx = ctx
	return b
}
