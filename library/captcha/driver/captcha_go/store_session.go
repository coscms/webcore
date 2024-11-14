package captcha_go

import (
	"context"

	"github.com/coscms/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewStoreSession() captcha.Storer {
	return &storeSession{}
}

const (
	sessionKey = `captchaGoGen`
)

type storeSession struct {
}

func (a *storeSession) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	eCtx.Session().Set(sessionKey, com.String(val))
	return nil
}

func (a *storeSession) Get(ctx context.Context, key string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	sessVal, ok := eCtx.Session().Get(sessionKey).(string)
	if !ok {
		return captcha.ErrIllegalKey
	}
	*(value.(*[]byte)) = com.Str2bytes(sessVal)
	return nil
}

func (a *storeSession) Delete(ctx context.Context, key string) error {
	eCtx := ctx.(echo.Context)
	eCtx.Session().Delete(sessionKey)
	return nil
}
