package pagebuilder

import (
	"github.com/webx-top/echo"
)

// NewPage 创建一个具有指定标题的新Page实例
func NewPage(title string) *Page {
	return &Page{Title: title}
}

// Page 页面
type Page struct {
	Template      string                     // 模板名称
	Title         string                     // 标题
	Header        string                     // 头部include html文件
	Data          interface{}                // 数据
	Body          echo.RenderContextWithData // 主体
	Footer        string                     // 尾部include html文件
	Breadcrumb    []*Breadcrumb              // 面包屑
	TopButtons    []*Button                  // 顶部按钮
	BottomButtons []*Button                  // 底部按钮
}

// SetData 设置页面数据并返回当前页面对象
// data: 要设置的页面数据
func (p *Page) SetData(data interface{}) *Page {
	p.Data = data
	return p
}

// SetBody 设置页面的主体内容并返回当前Page对象以便链式调用
func (p *Page) SetBody(body echo.RenderContextWithData) *Page {
	p.Body = body
	return p
}

// SetTemplate 设置页面使用的模板名称并返回当前Page对象以便链式调用
func (p *Page) SetTemplate(tmpl string) *Page {
	p.Template = tmpl
	return p
}

// SetHeader 设置页面的Header字段并返回Page对象以便链式调用
func (p *Page) SetHeader(header string) *Page {
	p.Header = header
	return p
}

// SetFooter 设置页面的页脚HTML文件并返回当前Page对象以便链式调用
func (p *Page) SetFooter(footer string) *Page {
	p.Footer = footer
	return p
}

// SetBreadcrumb 设置页面的面包屑导航
// 参数 breadcrumb 是面包屑导航项的切片
// 返回当前 Page 对象以便链式调用
func (p *Page) SetBreadcrumb(breadcrumb []*Breadcrumb) *Page {
	p.Breadcrumb = breadcrumb
	return p
}

// AddBreadcrumb 向页面添加面包屑导航项，支持添加多个面包屑项
// 返回当前Page对象以便链式调用
func (p *Page) AddBreadcrumb(breadcrumb ...*Breadcrumb) *Page {
	p.Breadcrumb = append(p.Breadcrumb, breadcrumb...)
	return p
}

// SetTopButtons 设置页面顶部的按钮列表并返回当前页面对象
func (p *Page) SetTopButtons(buttons []*Button) *Page {
	p.TopButtons = buttons
	return p
}

// AddTopButton 向页面顶部添加一个或多个按钮，并返回页面对象以便链式调用
func (p *Page) AddTopButton(buttons ...*Button) *Page {
	p.TopButtons = append(p.TopButtons, buttons...)
	return p
}

// SetBottomButtons 设置页面底部按钮列表并返回当前页面对象
func (p *Page) SetBottomButtons(buttons []*Button) *Page {
	p.BottomButtons = buttons
	return p
}

// AddBottomButton 向页面底部添加一个或多个按钮，并返回页面对象以便链式调用
func (p *Page) AddBottomButton(buttons ...*Button) *Page {
	p.BottomButtons = append(p.BottomButtons, buttons...)
	return p
}

// Render 渲染页面内容到响应中
// 根据页面数据类型自动选择默认模板（表格数据使用common/page_table，表单数据使用common/page_form）
// 参数:
//
//	ctx: echo框架的上下文对象
//
// 返回值:
//
//	error: 渲染过程中可能发生的错误
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
