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

func (a *storeCookie) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	var uid int64
	value := fmt.Sprintf(`%d|%d|%s`, time.Now().Unix(), uid, val)
	log.Debug(`set cookie value for CaptchaGo: `, value)
	eCtx.Cookie().EncryptSet(`CaptchaGo`, value, captcha.MaxAge)
	return nil
}

func (a *storeCookie) Get(ctx context.Context, key string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	cookieVal := eCtx.Cookie().DecryptGet(`CaptchaGo`)
	if len(cookieVal) == 0 {
		log.Debug(`failed to get cookie for CaptchaGo`)
		return captcha.ErrIllegalKey
	}
	parts := strings.SplitN(cookieVal, `|`, 3)
	if len(parts) != 3 {
		log.Debug(`failed to parse cookie for CaptchaGo: `, cookieVal)
		return captcha.ErrIllegalKey
	}
	if time.Now().Unix()-param.AsInt64(parts[0]) > captcha.MaxAge {
		log.Debug(`cookie has expired: CaptchaGo`)
		return captcha.ErrIllegalKey
	}
	*(value.(*[]byte)) = []byte(parts[2])
	return nil
}

func (a *storeCookie) Delete(ctx context.Context, key string) error {
	eCtx := ctx.(echo.Context)
	eCtx.Cookie().Set(`CaptchaGo`, ``, -1)
	return nil
}
