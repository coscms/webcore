package captcha_go

import (
	"context"

	"github.com/admpub/cache"
	"github.com/admpub/once"
	"github.com/coscms/captcha"
	"github.com/coscms/captcha/driver"
	captchaLib "github.com/coscms/webcore/library/captcha"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
}

var initialized once.Once

var DefaultStore captcha.Storer

func initialize() {
	var store captcha.Storer
	var err error
	if DefaultStore != nil {
		store = DefaultStore
	} else {
		store, err = cache.NewCacher(context.Background(), `memory`, cache.Options{Interval: int(captcha.MaxAge)})
		if err != nil {
			panic(err)
		}
	}
	err = driver.Initialize(store)
	if err != nil {
		panic(err)
	}
}
