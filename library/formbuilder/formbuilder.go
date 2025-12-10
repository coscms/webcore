package formbuilder

import (
	"errors"
	"strings"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	formsconfig "github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/formfilter"
	echoMw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
)

var (
	ErrJSONConfigFileNameInvalid = errors.New("*.form.json name invalid")
)

var LanguageConfigGetter func(echo.Context) language.Config
var Translateable bool

// New 创建表单构建器实例
func New(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	f := &FormBuilder{
		Forms:               forms.NewForms(forms.New()),
		hooks:               MethodHooks{},
		ctx:                 ctx,
		dbi:                 factory.DefaultDBI,
		langsGetter:         LanguageConfigGetter,
		formInputNamePrefix: FormInputNamePrefixDefault,
		translateable:       Translateable,
	}
	f.setDefaultLanguage()
	defaultHooks := []MethodHook{
		BindModel(f),
		ValidModel(f),
	}
	f.OnPost(defaultHooks...)
	f.OnPut(defaultHooks...)
	f.SetModel(model)
	f.Theme = common.BOOTSTRAP
	for _, option := range options {
		if option == nil {
			continue
		}
		option(f)
	}
	err := f.InitConfig()
	if err != nil {
		f.err = err
		f.exit = true
	}
	f.SetLabelFunc(func(txt string) string {
		return ctx.T(txt)
	})
	f.AddBeforeRender(func() {
		nextURL := ctx.Form(echo.DefaultNextURLVarName)
		if len(nextURL) > 0 {
			f.Elements(fields.HiddenField(echo.DefaultNextURLVarName).SetValue(nextURL))
		}
	})
	if csrfToken := ctx.Internal().String(echoMw.DefaultCSRFConfig.ContextKey); len(csrfToken) > 0 {
		f.AddBeforeRender(func() {
			f.Elements(fields.HiddenField(echoMw.DefaultCSRFConfig.ContextKey).SetValue(csrfToken))
		})
	}
	ctx.Set(`forms`, f.Forms)
	return f
}

// FormBuilder HTML表单构建器
type FormBuilder struct {
	*forms.Forms
	hooks               MethodHooks
	exit                bool
	err                 error
	ctx                 echo.Context
	configFile          string
	configPrepare       func(*formsconfig.Config) error
	config              *formsconfig.Config
	dbi                 *factory.DBI
	defaults            map[string]string
	filters             []formfilter.Options
	langsGetter         func(echo.Context) language.Config
	langDefault         string
	langConfig          *language.Config
	allowedNames        []string
	formInputNamePrefix string
	translateable       bool
}

// Exited 是否需要退出后续处理。此时一般有err值，用于记录错误原因
func (f *FormBuilder) Exited() bool {
	return f.exit
}

// Exit 设置退出标记
func (f *FormBuilder) Exit(exit ...bool) *FormBuilder {
	if len(exit) > 0 && !exit[0] {
		f.exit = false
	} else {
		f.exit = true
	}
	return f
}

// SetError 记录错误
func (f *FormBuilder) SetError(err error) *FormBuilder {
	f.err = err
	return f
}

// HasError 是否有错误
func (f *FormBuilder) HasError() bool {
	return f.err != nil
}

// Error 返回错误值
func (f *FormBuilder) Error() error {
	return f.err
}

// SetTranslateable sets whether the form fields should be translated.
// Returns the FormBuilder instance for method chaining.
func (f *FormBuilder) SetTranslateable(translateable bool) *FormBuilder {
	f.translateable = translateable
	return f
}

// RecvSubmission 接收客户端的提交
func (f *FormBuilder) RecvSubmission() error {
	ctx := f.ctx
	method := strings.ToUpper(ctx.Method())
	if f.err = f.hooks.Fire(method); f.err != nil {
		return f.err
	}
	f.err = f.hooks.Fire(`*`)
	if ctx.Response().Committed() {
		f.exit = true
	}
	return f.err
}

// Generate processes the form configuration, sets default values, and returns the FormBuilder instance for method chaining.
func (f *FormBuilder) Generate() *FormBuilder {
	f.ParseFromConfig()
	f.setDefaultValue()
	return f
}

// Snippet 表单片段
func (f *FormBuilder) Snippet() *FormBuilder {
	f.Config().Template = `allfields`
	f.Config().WithButtons = false
	return f
}

// FormData retrieves form data from the request based on the content type.
// Returns engine.URLValuer containing either POST form data or URL query parameters.
func (f *FormBuilder) FormData() engine.URLValuer {
	return FormData(f.ctx)
}
