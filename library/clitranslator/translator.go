package clitranslator

import (
	"github.com/admpub/once"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo/middleware/language"
)

var translate *language.Translate
var translock once.Once
var LangCode = `zh-CN` // 默认语言

func initTranslate() {
	translate = BuildTranslator(config.FromFile().Language, LangCode)
}

func GetTranslator() *language.Translate {
	translock.Do(initTranslate)
	return translate
}

func ResetTranslator() {
	translock.Reset()
}

func NewLanguage(c language.Config) *language.Language {
	c.SetFSFunc(bootconfig.LangFSFunc)
	c.Reload = false
	return language.New(&c)
}

func BuildTranslator(c language.Config, langCode string) *language.Translate {
	tr := &language.Translate{}
	lng := NewLanguage(c)
	tr.Reset(langCode, lng)
	return tr
}
