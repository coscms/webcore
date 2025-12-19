package formbuilder

import (
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/fields"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

// ClearCache 清除表单配置和模板缓存
func ClearCache() {
	common.ClearCachedConfig()
	common.ClearCachedTemplate()
}

// DelCachedConfig 删除指定的表单配置缓存
func DelCachedConfig(file string) bool {
	return common.DelCachedConfig(file)
}

// AddChoiceByKV adds choices to a form field based on key-value data.
func AddChoiceByKV(field fields.FieldInterface, kvData *echo.KVData, checkedKeys ...string) fields.FieldInterface {
	for _, kv := range kvData.Slice() {
		var checked bool
		if kv.H != nil {
			checked = kv.H.Bool(`checked`) || kv.H.Bool(`selected`)
		}
		if len(checkedKeys) > 0 {
			checked = com.InSlice(kv.K, checkedKeys)
		}
		field.AddChoice(kv.K, kv.V, checked)
	}
	return field
}

// SetChoiceByKV sets the choices of a form field based on key-value data.
func SetChoiceByKV(field fields.FieldInterface, kvData *echo.KVData, checkedKeys ...string) fields.FieldInterface {
	choices := []fields.InputChoice{}
	if len(checkedKeys) == 0 {
		switch f := field.(type) {
		case *fields.Field:
			if len(f.Value) > 0 {
				checkedKeys = append(checkedKeys, f.Value)
			}
		}
	}
	for _, kv := range kvData.Slice() {
		var checked bool
		if kv.H != nil {
			checked = kv.H.Bool(`checked`) || kv.H.Bool(`selected`)
		}
		if len(checkedKeys) > 0 {
			checked = com.InSlice(kv.K, checkedKeys)
		}
		choices = append(choices, fields.InputChoice{
			ID:      kv.K,
			Val:     kv.V,
			Checked: checked,
		})
	}

	field.SetChoices(choices)
	return field
}

// FormData retrieves form data from the HTTP request context.
// It automatically handles both application/x-www-form-urlencoded and multipart/form-data content types.
// Returns an URLValuer interface containing the parsed form data.
func FormData(ctx echo.Context) engine.URLValuer {
	contentType := ctx.Request().Header().Get(echo.HeaderContentType)
	var formData engine.URLValuer
	if strings.HasPrefix(contentType, echo.MIMEApplicationForm) {
		formData = ctx.Request().PostForm()
	} else {
		formData = ctx.Request().Form()
	}
	return formData
}
