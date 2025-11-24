package formbuilder

import "github.com/webx-top/echo/engine"

// langInputNamePrefix returns the input name prefix for the specified language in the format "Language[lang]"
func (f *FormBuilder) langInputNamePrefix(lang string) string {
	return `Language[` + lang + `]`
}

// SetLangInput sets the form input value for a specific language and field.
// Returns the FormBuilder instance for method chaining.
func (f *FormBuilder) SetLangInput(lang string, field string, value string, postFormOnly ...bool) *FormBuilder {
	if len(postFormOnly) > 0 && postFormOnly[0] {
		f.ctx.Request().PostForm().Set(f.langInputNamePrefix(lang)+`[`+field+`]`, value)
	} else {
		f.ctx.Request().Form().Set(f.langInputNamePrefix(lang)+`[`+field+`]`, value)
	}
	return f
}

// AnyLangInputCallback processes form input values for all languages, applying the given callback function to each language-specific input.
// The callback receives the current value and language code, and returns the modified value.
// If postFormOnly is true, only processes POST form data; otherwise processes both POST and GET data.
// Returns the FormBuilder instance for method chaining.
func (f *FormBuilder) AnyLangInputCallback(field string, callback func(value string, lang string) string, postFormOnly ...bool) *FormBuilder {
	var formData engine.URLValuer
	if len(postFormOnly) > 0 && postFormOnly[0] {
		formData = f.ctx.Request().PostForm()
	} else {
		formData = f.ctx.Request().Form()
	}
	for _, lang := range f.Languages().AllList {
		inputName := f.langInputNamePrefix(lang) + `[` + field + `]`
		value := formData.Get(inputName)
		value = callback(value, lang)
		formData.Set(inputName, value)
	}
	return f
}
