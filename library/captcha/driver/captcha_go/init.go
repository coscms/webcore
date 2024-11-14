package captcha_go

import (
	"encoding/gob"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/coscms/captcha"
	"github.com/coscms/captcha/driver"
	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/config"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
	func() {
		defer recover()
		gob.Register(map[string]int64{})
	}()
}

var initialized once.Once
var store captcha.Storer

func GetStore() captcha.Storer {
	initialized.Do(initialize)
	return store
}

var DefaultStore captcha.Storer = NewStoreSession()

func initialize() {
	cfg := config.FromFile().Extend.GetStore(`captchaGo`)
	switch cfg.String(`store`) {
	case `session`:
		if v, y := DefaultStore.(*storeSession); y {
			store = v
		} else {
			store = NewStoreSession()
		}
	case `cookie`:
		if v, y := DefaultStore.(*storeCookie); y {
			store = v
		} else {
			store = NewStoreCookie()
		}
	case `api`:
		if v, y := DefaultStore.(*storeAPI); y {
			store = v
		} else {
			apiURL := cfg.String(`apiURL`)
			if len(apiURL) == 0 {
				log.Warn(`captcha service apiURL is not set, default storage method is used`)
				store = DefaultStore
			} else {
				store = NewStoreAPI(apiURL, cfg.String(`secret`))
			}
		}
	default:
		store = DefaultStore
	}
	err := driver.Initialize(store)
	if err != nil {
		panic(err)
	}
}
