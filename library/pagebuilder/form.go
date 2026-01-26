package pagebuilder

import (
	"html/template"

	"github.com/coscms/webcore/library/formbuilder"
	"github.com/webx-top/echo"
)

// NewForm 创建一个新的表单实例，接收 echo.Context 和可选的 formbuilder.Option 参数
// 返回初始化后的 *Form 对象
func NewForm(ctx echo.Context, options ...formbuilder.Option) *Form {
	return &Form{options: options, FormBuilder: formbuilder.New(ctx, nil)}
}

type Form struct {
	options []formbuilder.Option
	*formbuilder.FormBuilder
}

func (f *Form) RenderContextWithData(ctx echo.Context, data interface{}) template.HTML {
	f.Init(data, f.options...)
	return f.FormBuilder.Render()
}
