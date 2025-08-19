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

func BuildTranslator(c language.Config, langCode string) *language.Translate {
	c.SetFSFunc(bootconfig.LangFSFunc)
	c.Reload = false
	lng := language.New(&c)
	tr := &language.Translate{}
	tr.Reset(langCode, lng)
	return tr
}
