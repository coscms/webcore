package formbuilder

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	formsconfig "github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"gopkg.in/yaml.v3"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	echoMw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/param"
)

var (
	ErrJSONConfigFileNameInvalid = errors.New("*.form.json name invalid")
)

var LanguagesDefaultGetter func(echo.Context) language.Config

// New 创建表单构建器实例
func New(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	f := &FormBuilder{
		Forms:       forms.NewForms(forms.New()),
		on:          MethodHooks{},
		ctx:         ctx,
		dbi:         factory.DefaultDBI,
		langsGetter: LanguagesDefaultGetter,
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
	on           MethodHooks
	exit         bool
	err          error
	ctx          echo.Context
	configFile   string
	config       *formsconfig.Config
	dbi          *factory.DBI
	defaults     map[string]string
	filters      []formfilter.Options
	langsGetter  func(echo.Context) language.Config
	langDefault  string
	langConfig   *language.Config
	allowedNames []string
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

// ParseConfigFile 解析配置文件 xxx.form.json
func (f *FormBuilder) ParseConfigFile(jsonformat ...bool) (*formsconfig.Config, error) {
	configFile := f.configFile + `.form`
	renderer, ok := f.ctx.Renderer().(driver.Driver)
	if !ok {
		return nil, fmt.Errorf(`FormBuilder: Expected renderer is "driver.Driver", but got "%T"`, f.ctx.Renderer())
	}
	var isJSON bool
	if len(jsonformat) > 0 {
		isJSON = jsonformat[0]
	}
	if isJSON {
		configFile += `.json`
	} else {
		configFile += `.yaml`
	}
	configFile = renderer.TmplPath(f.ctx, configFile)
	if len(configFile) == 0 {
		return nil, ErrJSONConfigFileNameInvalid
	}
	var cfg *formsconfig.Config
	b, err := renderer.RawContent(configFile)
	if err != nil || len(b) == 0 {
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf(`read file %s: %w`, configFile, err)
		}
		if renderer.Manager() == nil {
			return nil, fmt.Errorf(`renderer.Manager() is nil: %s`, configFile)
		}
		cfg = f.ToConfig()
		if isJSON {
			b, err = f.ToJSONBlob(cfg)
			if err != nil {
				return nil, fmt.Errorf(`[form.ToJSONBlob] %s: %w`, configFile, err)
			}
		} else {
			b, err = yaml.Marshal(cfg)
			if err != nil {
				return nil, fmt.Errorf(`[form:yaml.Marshal] %s: %w`, configFile, err)
			}
		}
		err = renderer.Manager().SetTemplate(configFile, b)
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, configFile, err)
		}
		f.ctx.Logger().Infof(f.ctx.T(`生成表单配置文件“%v”成功。`), configFile)
	} else {
		if isJSON {
			cfg, err = forms.Unmarshal(b, configFile)
			if err != nil {
				return nil, fmt.Errorf(`[forms.Unmarshal] %s: %w`, configFile, err)
			}
		} else {
			cfg, err = common.GetOrSetCachedConfig(configFile, func() (*formsconfig.Config, error) {
				cfg := &formsconfig.Config{}
				err := yaml.Unmarshal(b, cfg)
				return cfg, err
			})
			if err != nil {
				return nil, fmt.Errorf(`[form:yaml.Unmarshal] %s: %w`, configFile, err)
			}
		}
	}
	if cfg != nil {
		return cfg.Clone(), err
	}
	cfg = f.NewConfig()
	return cfg, err
}

func (f *FormBuilder) SetConfig(cfg *formsconfig.Config) *FormBuilder {
	f.config = cfg
	return f
}

func (f *FormBuilder) InitConfig() error {
	var cfg *formsconfig.Config
	var err error
	if f.config == nil {
		cfg, err = f.ParseConfigFile()
		if err != nil {
			return err
		}
	} else {
		cfg = f.config.Clone()
	}

	if f.Languages() != nil {
		f.toLangset(cfg)
	}

	f.Init(cfg)
	f.ParseFromConfig()

	defaultValues := f.DefaultValues()
	if len(defaultValues) > 0 {
		cfg.SetDefaultValue(func(fieldName string) string {
			val, ok := defaultValues[com.Title(fieldName)]
			if ok {
				return val
			}
			val = f.ctx.Form(fieldName)
			if len(val) == 0 && len(f.langDefault) > 0 {
				if after, found := strings.CutPrefix(fieldName, `Language[`+f.langDefault+`]`); found && len(after) > 0 {
					fieldName = strings.Trim(after, `[]`)
					val, _ = defaultValues[com.Title(fieldName)]
				}
			}
			return val
		})
	}
	return err
}

func (f *FormBuilder) toLangset(cfg *formsconfig.Config) {
	lgs := f.Languages()
	if lgs == nil {
		return
	}
	langCodes := lgs.AllList
	if len(langCodes) <= 1 {
		return
	}
	m, ok := f.Model.(factory.Short)
	if !ok {
		log.Warnf(`[formbuilder.toLangset] model %T does not implement factory.Short`, f.Model)
		return
	}
	var fields []string
	for _, info := range f.dbi.Fields[m.Short_()] {
		if info.Multilingual {
			fields = append(fields, info.GoName)
		}
	}
	if len(fields) == 0 {
		return
	}
	var setElems func(elems []*formsconfig.Element) []*formsconfig.Element
	setElems = func(elems []*formsconfig.Element) []*formsconfig.Element {
		var lastLangset *formsconfig.Element
		var lastLangsetIndex int
		var deleteIndexes []int
		for index, elem := range elems {
			if elem.Type == `fieldset` {
				elem.Elements = setElems(elem.Elements)
				elems[index] = elem
				continue
			}
			if elem.Type == `langset` {
				// 已经是langset类型，无需处理
				continue
			}
			if elem.Name == `` {
				continue
			}
			fieldName := com.Title(elem.Name)
			if !slices.Contains(fields, fieldName) {
				continue
			}
			cloned := elem.Clone()
			if lastLangset != nil && lastLangsetIndex == index-1 {
				// 紧跟在上一个langset后面，合并到上一个langset中
				lastLangset.Elements = append(lastLangset.Elements, cloned)
				deleteIndexes = append(deleteIndexes, index)
				lastLangsetIndex = index
				continue
			}
			// 创建新的langset
			elem.Type = `langset`
			elem.Elements = []*formsconfig.Element{cloned}
			for _, lang := range langCodes {
				label := lgs.ExtraBy(lang).String(`label`)
				if len(label) == 0 {
					label = lang
				}
				elem.AddLanguage(formsconfig.NewLanguage(lang, label, `Language[`+lang+`][%s]`))
			}
			lastLangset = elem
			lastLangsetIndex = index
		}
		// 删除已合并的元素
		if len(deleteIndexes) > 0 {
			newElems := []*formsconfig.Element{}
			for index, elem := range elems {
				if slices.Contains(deleteIndexes, index) {
					continue
				}
				newElems = append(newElems, elem)
			}
			elems = newElems
		}
		return elems
	}
	cfg.Elements = setElems(cfg.Elements)
}

// DefaultValues 获取model结构体各个字段在数据库中的默认值
func (f *FormBuilder) DefaultValues() map[string]string {
	if f.defaults != nil {
		return f.defaults
	}
	if f.dbi == nil || f.dbi.Fields == nil {
		return nil
	}
	m, ok := f.Model.(factory.Model)
	if !ok {
		return nil
	}
	fields, ok := f.dbi.Fields[m.Short_()]
	if !ok || fields == nil {
		return nil
	}
	f.defaults = map[string]string{}
	for _, info := range fields {
		v := m.GetField(info.GoName)
		var valStr string
		if v != nil && !reflect.ValueOf(v).IsZero() {
			valStr = param.AsString(v)
		}
		if len(valStr) == 0 {
			valStr = info.DefaultValue
		}
		f.defaults[info.GoName] = valStr
	}
	return f.defaults
}

// DefaultValue 查询某个结构体字段在数据库中对应的默认值
func (f *FormBuilder) DefaultValue(fieldName string) string {
	defaultValues := f.DefaultValues()
	if defaultValues == nil {
		return ``
	}
	fieldName = com.Title(fieldName)
	val, _ := defaultValues[fieldName]
	return val
}

// RecvSubmission 接收客户端的提交
func (f *FormBuilder) RecvSubmission() error {
	ctx := f.ctx
	method := strings.ToUpper(ctx.Method())
	if f.err = f.on.Fire(method); f.err != nil {
		return f.err
	}
	f.err = f.on.Fire(`*`)
	if ctx.Response().Committed() {
		f.exit = true
	}
	return f.err
}

// Snippet 表单片段
func (f *FormBuilder) Snippet() *FormBuilder {
	f.Config().Template = `allfields`
	return f
}

func (f *FormBuilder) setDefaultLanguage(langDefault ...string) *FormBuilder {
	var _langDefault string
	if len(langDefault) > 0 {
		_langDefault = langDefault[0]
	}
	if len(_langDefault) == 0 && f.Languages() != nil {
		_langDefault = f.Languages().Default
	}
	f.langDefault = _langDefault
	return f
}

func (f *FormBuilder) Languages() *language.Config {
	if f.langConfig != nil {
		return f.langConfig
	}
	if f.langsGetter != nil {
		c := f.langsGetter(f.ctx)
		f.langConfig = &c
		return f.langConfig
	}
	return nil
}
