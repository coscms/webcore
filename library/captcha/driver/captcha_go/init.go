package captcha_go

import (
	"context"
	"encoding/gob"

	"github.com/admpub/cache"
	_ "github.com/admpub/cache/redis5"
	"github.com/admpub/once"
	"github.com/coscms/captcha"
	"github.com/coscms/captcha/driver"
	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
)

func init() {
	captchaLib.Register(captchaLib.TypeGo, newCaptchaGo)
	func() {
		defer recover()
		gob.Register(map[string]int64{})
	}()
}

var initialized once.Once

var DefaultStore captcha.Storer

func initialize() {
	var store captcha.Storer
	var err error
	if DefaultStore != nil {
		store = DefaultStore
	} else {
		switch config.FromFile().Sys.SessionEngine {
		case `redis`:
			sessionConfig := config.FromFile().Sys.SessionConfig
			if sessionConfig == nil {
				sessionConfig = echo.H{}
			}
			store, err = cache.NewCacher(context.Background(), `redis`, cache.Options{
				Adapter:       `redis`,
				AdapterConfig: `network=` + sessionConfig.String(`network`, `tcp`) + `,addr=` + sessionConfig.String(`address`, `127.0.0.1:6379`) + `,password=` + sessionConfig.String(`password`) + `,db=` + sessionConfig.String(`db`, `2`) + `,pool_size=100,idle_timeout=180,hset_name=CaptchaGo,prefix=captchago:`,
				Interval:      int(captcha.MaxAge),
			})
		default:
			store, err = cache.NewCacher(context.Background(), `memory`, cache.Options{Interval: int(captcha.MaxAge)})
		}
		if err != nil {
			panic(err)
		}
	}
	err = driver.Initialize(store)
	if err != nil {
		panic(err)
	}
}
