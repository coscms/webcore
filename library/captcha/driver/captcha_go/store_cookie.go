package captcha_go

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewStoreCookie() captcha.Storer {
	return &storeCookie{}
}

type storeCookie struct {
}

const (
	cookieKey = `CaptchaGo`
)

func (a *storeCookie) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	var uid int64
	value := fmt.Sprintf(`%d|%d|%s`, time.Now().Unix(), uid, val)
	log.Debugf(`set cookie value for %s: %s`, cookieKey, value)
	eCtx.Cookie().EncryptSet(cookieKey, value, captcha.MaxAge)
	return nil
}

func (a *storeCookie) Get(ctx context.Context, key string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	cookieVal := eCtx.Cookie().DecryptGet(cookieKey)
	if len(cookieVal) == 0 {
		log.Debugf(`failed to get cookie for %s`, cookieKey)
		return captcha.ErrIllegalKey
	}
	parts := strings.SplitN(cookieVal, `|`, 3)
	if len(parts) != 3 {
		log.Debugf(`failed to parse cookie for %s: %s`, cookieKey, cookieVal)
		return captcha.ErrIllegalKey
	}
	if time.Now().Unix()-param.AsInt64(parts[0]) > captcha.MaxAge {
		log.Debugf(`cookie has expired: %s`, cookieKey)
		return captcha.ErrIllegalKey
	}
	*(value.(*[]byte)) = []byte(parts[2])
	return nil
}

func (a *storeCookie) Delete(ctx context.Context, key string) error {
	eCtx := ctx.(echo.Context)
	eCtx.Cookie().Set(cookieKey, ``, -1)
	return nil
}
