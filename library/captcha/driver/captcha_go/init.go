package captcha_go

import (
	"encoding/gob"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
	"github.com/coscms/webcore/cmd/bootconfig"
	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo/param"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
	bootconfig.OnStart(-1, func() {
		config.OnKeySetSettings(`captcha.go`, func(d config.Diff) error {
			if !d.IsDiff || d.Old == nil {
				return nil
			}
			cfg := param.AsStore(d.Old)
			cDriver := cfg.String(`driver`)
			if len(cDriver) > 0 {
				cType := cfg.GetStore(cDriver).String(`type`)
				if len(cType) > 0 {
					captcha.UnregisterInstance(cDriver, cType)
				}
			}
			return nil
		})
		config.OnKeySetSettings(`captcha.type`, func(d config.Diff) error {
			if !d.IsDiff || d.Old == nil {
				return nil
			}
			cType := param.AsString(d.Old)
			if cType == `go` {
				cfg := config.FromDB().GetStore(`captcha`)
				cfgo := cfg.GetStore(`go`)
				cDriver := cfgo.String(`driver`)
				if len(cDriver) > 0 {
					cType := cfgo.GetStore(cDriver).String(`type`)
					if len(cType) > 0 {
						captcha.UnregisterInstance(cDriver, cType)
					}
				}
			}
			return nil
		})
	})
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
