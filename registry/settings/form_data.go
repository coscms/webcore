package settings

import (
	"fmt"
	"reflect"

	"github.com/admpub/map2struct"
	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type DataDecoder func(v *dbschema.NgingConfig) (pointer interface{}, err error)

func (d DataDecoder) Register(name string) {
	RegisterDecoder(name, func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData, err := d(v)
		if err != nil {
			return err
		}
		if len(v.Value) > 0 {
			err = com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		if sd, ok := jsonData.(SetDefaults); ok {
			sd.SetDefaults()
		}
		r[`ValueObject`] = jsonData
		return err
	})
}

type DataDecoders map[string]DataDecoder

func (d DataDecoders) Register(group string) {
	for name, initor := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		initor.Register(name)
	}
}

type DataEncoder func(v *dbschema.NgingConfig, r echo.H) (pointer interface{}, err error)

func (d DataEncoder) Register(name string) {
	RegisterEncoder(name, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		cfg, err := d(v, r)
		if err != nil {
			return nil, err
		}
		if vd, ok := cfg.(Validator); ok {
			err = vd.Validate(v.Context())
			if err != nil {
				return nil, err
			}
		}
		return com.JSONEncode(cfg)
	})
}

type DataEncoders map[string]DataEncoder

func (d DataEncoders) Register(group string) {
	for name, from := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		from.Register(name)
	}
}

type SetDefaults interface {
	SetDefaults()
}

type Validator interface {
	Validate(echo.Context) error
}

func MakeDecoder(t reflect.Type) func(v *dbschema.NgingConfig, r echo.H) error {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf(`non-struct type is unsupported: %s`, t.Kind().String()))
	}
	return func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData := reflect.New(t).Interface()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		if sd, ok := jsonData.(SetDefaults); ok {
			sd.SetDefaults()
		}
		r[`ValueObject`] = jsonData
		return nil
	}
}

func MakeEncoder(t reflect.Type) func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf(`non-struct type is unsupported: %s`, t.Kind().String()))
	}
	return func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		cfg := reflect.New(t).Interface()
		err := map2struct.Scan(cfg, r, `json`)
		if err != nil {
			return nil, err
		}
		if vd, ok := cfg.(Validator); ok {
			err = vd.Validate(v.Context())
			if err != nil {
				return nil, err
			}
		}
		return com.JSONEncode(cfg)
	}
}
