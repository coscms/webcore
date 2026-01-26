package formbuilder

import (
	"errors"
	"strings"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	formsconfig "github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"

	"github.com/webx-top/com"
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
var TranslateableGetter func(echo.Context) bool

func NewSnippet(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	options = append(options, Snippet(true))
	return New(ctx, model, options...)
}

// New 创建表单构建器实例
func New(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	f := &FormBuilder{
		Forms:               forms.NewForms(forms.New()),
		hooks:               MethodHooks{},
		ctx:                 ctx,
		dbi:                 factory.DefaultDBI,
		langsGetter:         LanguageConfigGetter,
		formInputNamePrefix: FormInputNamePrefixDefault,
		translateable:       TranslateableGetter,
		ctxStoreKey:         `forms`,
	}
	if model == nil && len(options) == 0 {
		return f
	}
	return f.Init(model, options...)
}

func (f *FormBuilder) Init(model interface{}, options ...Option) *FormBuilder {
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
		return f.ctx.T(txt)
	}).SetNameFunc(com.CamelCase)
	f.AddBeforeRender(func() {
		nextURL := f.ctx.Form(echo.DefaultNextURLVarName)
		if len(nextURL) > 0 {
			f.Elements(fields.HiddenField(echo.DefaultNextURLVarName).SetValue(nextURL))
		}
	})
	if csrfToken := f.ctx.Internal().String(echoMw.DefaultCSRFConfig.ContextKey); len(csrfToken) > 0 {
		f.AddBeforeRender(func() {
			f.Elements(fields.HiddenField(echoMw.DefaultCSRFConfig.ContextKey).SetValue(csrfToken))
		})
	}
	if len(f.ctxStoreKey) > 0 {
		f.ctx.Set(f.ctxStoreKey, f.Forms)
	}
	return f
}

// FormBuilder HTML表单构建器
type FormBuilder struct {
	*forms.Forms
	hooks               MethodHooks
	exit                bool
	snippet             bool
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
	translateable       func(echo.Context) bool
	ctxStoreKey         string
	translateLabelCols  int
	renames             map[string]string
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

// SetTranslateLabelCols 设置翻译标签的列数
func (f *FormBuilder) SetTranslateLabelCols(cols int) *FormBuilder {
	f.translateLabelCols = cols
	return f
}

// SetTranslateable sets the function that determines if the form should be translated.
// Returns the FormBuilder instance for method chaining.
func (f *FormBuilder) SetTranslateable(translateable func(echo.Context) bool) *FormBuilder {
	f.translateable = translateable
	return f
}

// Translateable returns whether the form builder supports translation capabilities.
// It checks if the translateable function is set and calls it with the current context.
// Returns false if no translateable function is configured.
func (f *FormBuilder) Translateable() bool {
	if f.translateable == nil {
		return false
	}
	return f.translateable(f.ctx)
}

// SetRenames 设置字段重命名
func (f *FormBuilder) SetRenames(renames map[string]string) *FormBuilder {
	f.renames = renames
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
func (f *FormBuilder) Generate(emptyValue ...bool) *FormBuilder {
	f.ParseFromConfig()
	if len(emptyValue) > 0 && emptyValue[0] {
		return f
	}
	f.setDefaultValue()
	return f
}

// SetSnippet sets whether to generate a snippet for the form.
// Returns the FormBuilder instance for method chaining.
// If snippet is true, the form will generate a snippet containing only the form fields.
// Otherwise, the form will generate a full HTML page.
func (f *FormBuilder) SetSnippet(snippet bool) *FormBuilder {
	f.snippet = snippet
	return f
}

// Snippet 表单片段
func (f *FormBuilder) Snippet() *FormBuilder {
	return f.setSnippetConfig(f.Config())
}

func (f *FormBuilder) setSnippetConfig(cfg *formsconfig.Config) *FormBuilder {
	if cfg.Template != `allfields` {
		cfg.Template = `allfields`
	}
	if cfg.WithButtons {
		cfg.WithButtons = false
	}
	return f
}

// FormData retrieves form data from the request based on the content type.
// Returns engine.URLValuer containing either POST form data or URL query parameters.
func (f *FormBuilder) FormData() engine.URLValuer {
	return FormData(f.ctx)
}
