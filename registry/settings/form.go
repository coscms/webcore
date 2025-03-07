package settings

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/admpub/map2struct"
	"github.com/coscms/forms"
	formsconfig "github.com/coscms/forms/config"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/errorslice"
	_ "github.com/coscms/webcore/library/formbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewForm(group string, short string, label string, opts ...FormSetter) *SettingForm {
	f := &SettingForm{
		Group: group,
		Short: short,
		Label: label,
		items: map[string]*dbschema.NgingConfig{},
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type SettingForm struct {
	Short        string   //简短标签
	Label        string   //标签文本
	Group        string   //组标识
	Tmpl         []string //输入表单模板路径
	HeadTmpl     []string
	FootTmpl     []string
	items        map[string]*dbschema.NgingConfig //配置项
	hookPost     []func(echo.Context) error       //数据提交逻辑处理
	hookGet      []func(echo.Context) error       //数据读取逻辑处理
	renderer     func(echo.Context) template.HTML
	config       *formsconfig.Config
	dataDecoders DataDecoders //Decoder: table => form
	dataEncoders DataEncoders //Encoder: form => table
}

func (s *SettingForm) Merge(source *SettingForm) *SettingForm {
	if len(source.HeadTmpl) > 0 {
		s.HeadTmpl = append(s.HeadTmpl, source.HeadTmpl...)
	}
	if len(source.Tmpl) > 0 {
		s.Tmpl = append(s.Tmpl, source.Tmpl...)
	}
	if len(source.FootTmpl) > 0 {
		s.FootTmpl = append(s.FootTmpl, source.FootTmpl...)
	}
	if len(source.Short) > 0 && len(s.Short) == 0 {
		s.Short = source.Short
	}
	if len(source.Label) > 0 && len(s.Label) == 0 {
		s.Label = source.Label
	}
	for k, v := range source.items {
		s.items[k] = v
	}
	if len(source.hookPost) > 0 {
		s.hookPost = append(s.hookPost, source.hookPost...)
	}
	if len(source.hookGet) > 0 {
		s.hookGet = append(s.hookGet, source.hookGet...)
	}
	if source.renderer != nil {
		if s.renderer == nil {
			s.renderer = source.renderer
		} else {
			renderer := s.renderer
			s.renderer = func(ctx echo.Context) template.HTML {
				return renderer(ctx) + source.renderer(ctx)
			}
		}
	}
	if source.config != nil {
		if s.config == nil {
			s.config = source.config
		} else {
			s.config.Merge(source.config)
		}
	}
	if source.dataDecoders != nil {
		if s.dataDecoders == nil {
			s.dataDecoders = source.dataDecoders
		} else {
			for k, v := range source.dataDecoders {
				s.dataDecoders[k] = v
			}
		}
	}
	if source.dataEncoders != nil {
		if s.dataEncoders == nil {
			s.dataEncoders = source.dataEncoders
		} else {
			for k, v := range source.dataEncoders {
				s.dataEncoders[k] = v
			}
		}
	}
	return s
}

func (s *SettingForm) AddTmpl(tmpl ...string) *SettingForm {
	s.Tmpl = append(s.Tmpl, tmpl...)
	return s
}

func (s *SettingForm) AddHeadTmpl(tmpl ...string) *SettingForm {
	s.HeadTmpl = append(s.HeadTmpl, tmpl...)
	return s
}

func (s *SettingForm) AddFootTmpl(tmpl ...string) *SettingForm {
	s.FootTmpl = append(s.FootTmpl, tmpl...)
	return s
}

func (s *SettingForm) AddHookPost(hook func(echo.Context) error) *SettingForm {
	s.hookPost = append(s.hookPost, hook)
	return s
}

func (s *SettingForm) SetFormConfig(formcfg *formsconfig.Config) *SettingForm {
	if len(formcfg.Theme) == 0 {
		formcfg.Theme = `bootstrap3`
	}
	s.config = formcfg
	return s
}

// SetDataTransfer 数据转换
// dataDecoder: Decoder(table => form)
// dataEncoder: Encoder(form => table)
func (s *SettingForm) SetDataTransfer(name string, dataDecoder DataDecoder, dataEncoder DataEncoder) *SettingForm {
	if dataDecoder != nil {
		if s.dataDecoders == nil {
			s.dataDecoders = DataDecoders{}
		}
		s.dataDecoders[name] = dataDecoder
	}
	if dataEncoder != nil {
		if s.dataEncoders == nil {
			s.dataEncoders = DataEncoders{}
		}
		s.dataEncoders[name] = dataEncoder
	}
	return s
}

func (s *SettingForm) SetTransferType(name string, dest interface{}) *SettingForm {
	rType := GetReflectType(dest)
	if rType.Kind() != reflect.Struct {
		panic(fmt.Sprintf(`non-struct type is unsupported: %s`, rType.Kind().String()))
	}
	return s.SetDataTransfer(name, func(v *dbschema.NgingConfig) (interface{}, error) {
		return reflect.New(rType).Interface(), nil
	}, func(v *dbschema.NgingConfig, r echo.H) (interface{}, error) {
		cfg := reflect.New(rType).Interface()
		err := map2struct.Scan(cfg, r, `json`)
		return cfg, err
	})
}

func (s *SettingForm) AddConfig(configs ...*dbschema.NgingConfig) *SettingForm {
	if s.items == nil {
		s.items = map[string]*dbschema.NgingConfig{}
	}
	for _, c := range configs {
		if c.Group != s.Group {
			c.Group = s.Group
		}
		s.items[c.Key] = c
	}
	return s
}

func (s *SettingForm) SetRenderer(renderer func(echo.Context) template.HTML) *SettingForm {
	s.renderer = renderer
	return s
}

func StructFieldConvert(s string) string {
	// group[key][value][objkey]
	p := strings.Index(s, `[`)
	if p < 0 {
		return s
	}
	s = s[p:]
	p = strings.Index(s, `[value]`)
	if p < 0 {
		return s
	}
	s = s[:p] + s[p+len(`[value]`):]
	return s
}

func FormStoreToMap(d echo.H) echo.H {
	m := echo.H{}
	for k, v := range d {
		// api.ValueObject
		vo := param.AsStore(v).Get(`ValueObject`)
		if vo != nil {
			m[k] = vo
		}
	}
	return m
}

func (s *SettingForm) Render(ctx echo.Context) template.HTML {
	if s.renderer != nil {
		return s.renderer(ctx)
	}
	if s.config != nil {
		form := forms.NewForms(forms.NewWithConfig(s.config))
		form.SetStructFieldConverter(StructFieldConvert)
		if d, y := ctx.Get(s.Group).(echo.H); y {
			m := FormStoreToMap(d)
			form.SetModel(m)
		}
		form.ParseFromConfig(true)
		return form.Render()
	}
	var htmlContent string
	var stored echo.Store
	if fn, ok := ctx.GetFunc(`Stored`).(func() echo.Store); ok {
		stored = fn()
	} else {
		stored = ctx.Stored()
	}
	for _, t := range s.Tmpl {
		if len(t) == 0 {
			continue
		}
		b, err := ctx.Fetch(t, stored)
		if err != nil {
			htmlContent += err.Error()
		} else {
			htmlContent += string(b)
		}
	}
	return template.HTML(htmlContent)
}

func (s *SettingForm) AddHookGet(hook func(echo.Context) error) *SettingForm {
	s.hookGet = append(s.hookGet, hook)
	return s
}

func (s *SettingForm) RunHookPost(ctx echo.Context) error {
	if s.hookPost == nil {
		return nil
	}
	errs := errorslice.New()
	for _, hook := range s.hookPost {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}

func (s *SettingForm) RunHookGet(ctx echo.Context) error {
	if s.hookGet == nil {
		return nil
	}
	errs := errorslice.New()
	for _, hook := range s.hookGet {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}
