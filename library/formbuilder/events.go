package formbuilder

import (
	"errors"
	"slices"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/validation"
)

type MethodHook func() error
type MethodHooks map[string][]MethodHook

func (hooks MethodHooks) On(method string, funcs ...MethodHook) {
	if _, ok := hooks[method]; !ok {
		hooks[method] = []MethodHook{}
	}
	hooks[method] = append(hooks[method], funcs...)
}

func (hooks MethodHooks) Off(methods ...string) {
	for _, method := range methods {
		delete(hooks, method)
	}
}

func (hooks MethodHooks) OffAll() {
	for method := range hooks {
		delete(hooks, method)
	}
}

func (hooks MethodHooks) Fire(method string) error {
	funcs, ok := hooks[method]
	if !ok {
		return nil
	}
	var err error
	for _, fn := range funcs {
		if err = fn(); err != nil {
			return err
		}
	}
	return err
}

func BindModel(form *FormBuilder) MethodHook {
	return func() error {
		names := form.Config().GetNames()
		if form.langDefault != "" && form.Languages() != nil {
			if form.ctx.Lang().Normalize() == form.langDefault {
				//langKey := com.UpperCaseFirst(form.langDefault)
				formData := form.FormData()
				for _, name := range names {
					//pp.Println(name)
					if after, found := strings.CutPrefix(name, form.langInputNamePrefix(form.langDefault)); found {
						nameRaw := strings.Trim(after, `[]`)
						names = append(names, nameRaw)
						nameLower := com.LowerCaseFirst(nameRaw)
						formName := form.langInputNamePrefix(form.langDefault) + `[` + nameLower + `]`
						values := form.ctx.FormValues(formName)
						if len(values) == 0 {
							formName = form.langInputNamePrefix(form.langDefault) + `[` + nameRaw + `]`
							values = form.ctx.FormValues(formName)
						}
						//pp.Println(formName)
						echo.SetFormValues(formData, nameLower, values)
					}
				}
			}
			//pp.Println(form.ctx.Lang().Normalize(), form.langDefault, form.ctx.Forms())
		}
		if len(form.allowedNames) > 0 {
			for _, name := range form.allowedNames {
				if !slices.Contains(names, name) {
					names = append(names, name)
				}
			}
		}
		opts := []formfilter.Options{formfilter.Include(names...)}
		opts = append(opts, form.filters...)
		return form.ctx.MustBind(form.Model, formfilter.Build(opts...))
	}
}

func ValidModel(form *FormBuilder) MethodHook {
	return func() error {
		form.ValidFromConfig()
		err := form.Validate().Error()
		if !errors.Is(err, validation.NoError) {
			form.ctx.Data().SetInfo(err.Message, 0).SetZone(err.Field)
			return err
		}
		return nil
	}
}
