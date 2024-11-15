package captcha_go

import (
	"encoding/gob"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
	"github.com/coscms/webcore/cmd/bootconfig"
	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
	"golang.org/x/sync/singleflight"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
	bootconfig.OnStart(-1, func() {
		config.OnKeySetSettings(`captcha.go`, func(d config.Diff) error {
			if !d.IsDiff || d.Old == nil {
				return nil
			}
			cfg := param.AsStore(d.Old)
			unregister(cfg)
			return nil
		})
		config.OnKeySetSettings(`captcha.type`, func(d config.Diff) error {
			if !d.IsDiff || d.Old == nil {
				return nil
			}
			cType := param.AsString(d.Old)
			if cType == `go` {
				cfg := config.FromDB().GetStore(`captcha`)
				cfgo := cfg.GetStore(cType)
				unregister(cfgo)
			}
			return nil
		})
	})
	func() {
		defer recover()
		gob.Register(map[string]int64{})
	}()
}

func unregister(cfgo echo.H) {
	cDriver := cfgo.String(`driver`)
	if len(cDriver) > 0 {
		cType := cfgo.GetStore(cDriver).String(`type`)
		if len(cType) > 0 {
			captcha.UnregisterInstance(cDriver, cType)
		}
	}
}

var DefaultStorer captcha.Storer = NewStoreSession()
var StorerConstructor = storerDefaultConstructor
var usedStorer atomic.Value
var sg singleflight.Group

func GetStorer() (captcha.Storer, error) {
	store, ok := usedStorer.Load().(captcha.Storer)
	if ok {
		return store, nil
	}
	val, err, _ := sg.Do(`initStorer`, func() (interface{}, error) {
		st, err := StorerConstructor()
		if err != nil {
			log.Errorf(`captchaGo initialization storage engine failed: %v`, err)
			return st, err
		}
		log.Okayf(`captchaGo uses storage engine: %T`, st)
		usedStorer.Store(st)
		return st, err
	})
	if err != nil {
		return nil, err
	}
	return val.(captcha.Storer), nil
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
