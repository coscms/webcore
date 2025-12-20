package formbuilder

import (
	"slices"
	"strings"

	"github.com/admpub/log"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/tplfunc"
)

// setDefaultLanguage sets the default language for the form builder.
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

// Languages returns the language configuration for the form builder.
// If a custom language getter is set, it will be called to retrieve the configuration.
// Returns nil if no language configuration is available.
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

// setMultilingualElems processes form elements to group multilingual fields into langset elements.
// It checks each element against the provided list of multilingual fields and restructures
// the elements accordingly. Elements that are part of a langset are grouped together,
// and new langset elements are created as needed.
func (f *FormBuilder) setMultilingualElems(multilingualFields []string, elems []*formsconfig.Element) []*formsconfig.Element {
	lgs := f.Languages()
	var lastLangset *formsconfig.Element
	var lastLangsetIndex int
	var deleteIndexes []int
	for index, elem := range elems {
		if elem.Type == `fieldset` {
			elem.Elements = f.setMultilingualElems(multilingualFields, elem.Elements)
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
		fieldName := elem.Name
		if strings.HasSuffix(fieldName, `]`) { // 处理数组字段名，如 "tags[name]"
			start := strings.LastIndex(fieldName, `[`)
			if start > -1 {
				fieldName = fieldName[start+1 : len(fieldName)-1]
			}
		}
		fieldName = com.Title(fieldName)
		if !slices.Contains(multilingualFields, fieldName) {
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
		if len(elem.Template) > 0 {
			elem.Template = ``
		}
		for _, lang := range lgs.AllList {
			label := lgs.ExtraBy(lang).String(`label`)
			if len(label) == 0 {
				label = lang
			}
			elem.AddLanguage(formsconfig.NewLanguage(lang, label, f.langInputNamePrefix(lang)+`[%s]`))
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

// toLangset converts multilingual form fields into langset elements in the form configuration.
// It processes the form elements recursively, grouping multilingual fields under langset containers.
// For each multilingual field, it creates language-specific inputs based on the available languages.
// Fields that implement factory.Short interface are checked for multilingual support.
// The function modifies the provided config in-place and does not return any value.
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
	var multilingualFields []string
	for _, info := range f.dbi.Fields[m.Short_()] {
		if info.Multilingual {
			multilingualFields = append(multilingualFields, info.GoName)
		}
	}
	if len(multilingualFields) == 0 {
		return
	}
	cfg.Elements = f.setMultilingualElems(multilingualFields, cfg.Elements)
	if f.Translateable() {
		cfg.Elements = append(cfg.Elements, &formsconfig.Element{
			ID:        f.ctxStoreKey + `ForceTranslate` + tplfunc.RandomString(6),
			Type:      `checkbox`,
			Name:      `forceTranslate`,
			Label:     f.ctx.T(`自动翻译`),
			LabelCols: f.translateLabelCols,
			Choices: []*formsconfig.Choice{
				{
					Option: []string{`1`, f.ctx.T(`是`)},
				},
			},
		})
	}
}
