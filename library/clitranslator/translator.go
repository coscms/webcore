package clitranslator

import (
	"github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo/middleware/language"
)

var translate *language.Translate
var translock once.Once
var LangCode = `zh-CN` // 默认语言

func initTranslate() {
	translate = config.FromFile().BuildTranslator(LangCode)
}

func GetTranslator() *language.Translate {
	translock.Do(initTranslate)
	return translate
}

func ResetTranslator() {
	translock.Reset()
}
