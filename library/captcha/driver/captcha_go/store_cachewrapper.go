package captcha_go

import (
	"context"

	"github.com/coscms/captcha"
	"github.com/webx-top/echo"
)

func NewStoreCacheWrapper(cacher captcha.Storer) captcha.Storer {
	return &storeCacheWrapper{cacher: cacher}
}

type storeCacheWrapper struct {
	cacher captcha.Storer
}

func (a *storeCacheWrapper) Put(ctx context.Context, _ string, val interface{}, timeout int64) error {
	eCtx := ctx.(echo.Context)
	a.cacher.Put(ctx, sessionKey+eCtx.Session().MustID(), val, timeout)
	return nil
}

func (a *storeCacheWrapper) Get(ctx context.Context, _ string, value interface{}) error {
	eCtx := ctx.(echo.Context)
	return a.cacher.Get(ctx, sessionKey+eCtx.Session().MustID(), value)
}

func (a *storeCacheWrapper) Delete(ctx context.Context, _ string) error {
	eCtx := ctx.(echo.Context)
	return a.cacher.Delete(ctx, sessionKey+eCtx.Session().MustID())
}
