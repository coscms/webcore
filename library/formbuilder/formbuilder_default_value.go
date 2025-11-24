package formbuilder

import (
	"reflect"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

// setDefaultValue sets default values for form fields based on DefaultValues map.
// It handles both direct field names and language-specific field names (with prefix).
// If no default values are provided, the function does nothing.
// The function also checks for form input values as fallback defaults.
func (f *FormBuilder) setDefaultValue() {
	defaultValues := f.DefaultValues()
	if len(defaultValues) == 0 {
		return
	}
	// 需要先调用 ParseFromConfig() 来生成多语言输入表单域
	f.Config().SetDefaultValue(func(fieldName string) string {
		val, ok := defaultValues[com.Title(fieldName)]
		if ok {
			return val
		}
		val = f.ctx.Form(fieldName)
		if len(val) == 0 && len(f.langDefault) > 0 {
			if after, found := strings.CutPrefix(fieldName, f.langInputNamePrefix(f.langDefault)); found && len(after) > 0 {
				fieldName = strings.Trim(after, `[]`)
				val, _ = defaultValues[com.Title(fieldName)]
			}
		}
		return val
	})
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
