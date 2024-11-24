package settings

import (
	"reflect"

	"github.com/coscms/webcore/dbschema"
)

func AddConfigs(configs map[string]map[string]*dbschema.NgingConfig) {
	for group, configs := range configs {
		AddDefaultConfig(group, configs)
	}
}

func AddTmpl(group string, tmpl string, opts ...FormSetter) {
	// 注册配置模板和逻辑
	index, setting := Get(group)
	if index == -1 || setting == nil {
		return
	}
	if len(tmpl) > 0 {
		setting.AddTmpl(tmpl)
	}
	for _, option := range opts {
		option(setting)
	}
}

// RegisterTransferType 注册数据转换类型
// 名称支持"group"或"group.key"两种格式，例如:
// settings.RegisterTransferType(`sms`,...)对整个sms组的配置有效
// settings.RegisterTransferType(`sms.twilio`,...)对sms组内key为twilio的配置有效
func RegisterTransferType(group string, dest interface{}) {
	rType := GetReflectType(dest)
	RegisterDecoder(group, MakeDecoder(rType))
	RegisterEncoder(group, MakeEncoder(rType))
}

func GetReflectType(dest interface{}) reflect.Type {
	var rType reflect.Type
	switch v := dest.(type) {
	case reflect.Type:
		rType = v
	case reflect.Value:
		v = reflect.Indirect(v)
		rType = v.Type()
	default:
		rValue := reflect.ValueOf(dest)
		rValue = reflect.Indirect(rValue)
		rType = rValue.Type()
	}
	return rType
}
