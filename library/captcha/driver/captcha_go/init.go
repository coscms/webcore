package captcha_go

import (
	"encoding/gob"

	"github.com/admpub/once"
	"github.com/coscms/captcha"
	"github.com/coscms/captcha/driver"
	captchaLib "github.com/coscms/webcore/library/captcha"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
	func() {
		defer recover()
		gob.Register(map[string]int64{})
	}()
}

var initialized once.Once

var DefaultStore captcha.Storer = NewStoreCookie()

func initialize() {
	err := driver.Initialize(DefaultStore)
	if err != nil {
		panic(err)
	}
}
