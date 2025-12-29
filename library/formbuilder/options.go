package formbuilder

import (
	"github.com/coscms/forms/config"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/middleware/language"
)

type Option func(*FormBuilder)

// IgnoreFields 疏略某些字段的验证
func IgnoreFields(ignoreFields ...string) Option {
	return func(form *FormBuilder) {
		form.CloseValid(ignoreFields...)
	}
}

// Theme 设置forms模板风格
func Theme(theme string) Option {
	return func(form *FormBuilder) {
		form.Theme = theme
	}
}

// FormFilter 设置表单过滤
func FormFilter(filters ...formfilter.Options) Option {
	return func(form *FormBuilder) {
		form.filters = filters
	}
}

// ConfigFile 指定要解析的配置文件。如果silent=true则仅仅设置配置文件而不解析
func ConfigFile(jsonFile string, silent ...bool) Option {
	return func(f *FormBuilder) {
		f.configFile = jsonFile
		if len(silent) > 0 && silent[0] {
			return
		}
		cfg, err := f.ParseConfigFile()
		if err != nil {
			panic(err)
		}
		f.SetConfig(cfg)
	}
}

// Config 指定表单配置
func Config(cfg *config.Config) Option {
	return func(f *FormBuilder) {
		f.SetConfig(cfg)
	}
}

func ConfigPrepare(configPrepare func(*formsconfig.Config) error) Option {
	return func(f *FormBuilder) {
		f.configPrepare = configPrepare
	}
}

// RenderBefore 设置渲染表单前的钩子函数
func RenderBefore(fn func()) Option {
	return func(f *FormBuilder) {
		f.AddBeforeRender(fn)
	}
}

// DBI 指定模型数据表所属DBI(数据库信息)
func DBI(dbi *factory.DBI) Option {
	return func(f *FormBuilder) {
		f.dbi = dbi
	}
}

// LanguagesGetter 多语言配置
func LanguagesGetter(langsGetter func(echo.Context) language.Config, langDefault ...string) Option {
	return func(f *FormBuilder) {
		f.langsGetter = langsGetter
		f.setDefaultLanguage(langDefault...)
	}
}

// AllowedNames sets the allowed field names for the form builder.
// It accepts a variadic list of string names to be permitted in the form.
func AllowedNames(names ...string) Option {
	return func(f *FormBuilder) {
		f.allowedNames = names
	}
}

// FormInputNamePrefix sets the prefix for form input names in the FormBuilder.
func FormInputNamePrefix(prefix string) Option {
	return func(f *FormBuilder) {
		f.formInputNamePrefix = prefix
	}
}

// Translateable sets the function that determines if the form should be translated.
// Returns the FormBuilder instance for method chaining.
func Translateable(translateable func(echo.Context) bool) Option {
	return func(f *FormBuilder) {
		f.translateable = translateable
	}
}

// ContextStoreKey sets the context store key for storing forms in the Echo context.
func ContextStoreKey(key string) Option {
	return func(f *FormBuilder) {
		f.ctxStoreKey = key
	}
}

// TranslateLabelCols 设置翻译标签的列数
func TranslateLabelCols(cols int) Option {
	return func(f *FormBuilder) {
		f.translateLabelCols = cols
	}
}

// Renames 设置字段重命名
func Renames(renames map[string]string) Option {
	return func(f *FormBuilder) {
		f.renames = renames
	}
}
