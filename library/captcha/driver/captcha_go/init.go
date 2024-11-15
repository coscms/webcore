package captcha_go

import (
	"encoding/gob"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
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

var DefaultStorer captcha.Storer = NewStoreSession()
var StorerConstructor = storerDefaultConstructor
var usedStorer atomic.Value

func GetStorer() (captcha.Storer, error) {
	store, ok := usedStorer.Load().(captcha.Storer)
	if ok {
		return store, nil
	}
	store, err := StorerConstructor()
	if err != nil {
		return nil, err
	}
	usedStorer.Store(store)
	return store, nil
}

func storerDefaultConstructor() (captcha.Storer, error) {
	cfg := config.FromFile().Extend.GetStore(`captchaGo`)
	var store captcha.Storer
	switch cfg.String(`store`) {
	case `session`:
		if v, y := DefaultStorer.(*storeSession); y {
			store = v
		} else {
			store = NewStoreSession()
		}
	case `cookie`:
		if v, y := DefaultStorer.(*storeCookie); y {
			store = v
		} else {
			store = NewStoreCookie()
		}
	case `api`:
		if v, y := DefaultStorer.(*storeAPI); y {
			store = v
		} else {
			apiURL := cfg.String(`apiURL`)
			if len(apiURL) == 0 {
				log.Warn(`captcha service apiURL is not set, default storage method is used`)
				store = DefaultStorer
			} else {
				store = NewStoreAPI(apiURL, cfg.String(`secret`))
			}
		}
	default:
		store = DefaultStorer
	}
	return store, nil
}
