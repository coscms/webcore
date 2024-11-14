package captcha_go

import (
	"time"

	"github.com/coscms/captcha"
	"github.com/webx-top/echo"
)

func RegisterRoute(g echo.RouteRegister) {
	g.Get(`/:driver/:type`, CaptchaGoData)
	g.Post(`/:driver/:type`, CaptchaGoVerify)
}

func CaptchaGoData(ctx echo.Context) error {
	resp := captcha.APIResponse{}
	c, err := captcha.GetInstance(ctx.Param(`driver`), ctx.Param(`type`))
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	data, err := c.MakeData(ctx)
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	return ctx.JSON(resp.SetData(data))
}

func CaptchaGoVerify(ctx echo.Context) error {
	resp := captcha.APIResponse{}
	c, err := captcha.GetInstance(ctx.Param(`driver`), ctx.Param(`type`))
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	key := ctx.Form(`key`)
	err = c.Verify(ctx, key, ctx.Form(`response`))
	if err != nil {
		return ctx.JSON(resp.SetError(err.Error()))
	}
	captchaGoSetSuccessKey(ctx, key)
	c.Storer().Delete(ctx, key)
	return ctx.JSON(resp.SetSuccess())
}

const captchaGoSessionKey = `captchaGoKey`

func captchaGoSetSuccessKey(ctx echo.Context, key string) {
	secrets, ok := ctx.Session().Get(captchaGoSessionKey).(map[string]int64)
	if !ok {
		secrets = map[string]int64{}
	} else {
		// 始终只保存10条，避免session数据无限增长
		if len(secrets) >= 10 {
			var minTs int64
			var minKey string
			for k, v := range secrets {
				if minTs == 0 || minTs > v {
					minTs = v
					minKey = k
				}
			}
			delete(secrets, minKey)
		}
	}
	secrets[key] = time.Now().Unix()
	ctx.Session().Set(captchaGoSessionKey, secrets)
}

func captchaGoVerifySuccessKey(ctx echo.Context, key string, delAfterVerfiy bool) bool {
	secrets, ok := ctx.Session().Get(captchaGoSessionKey).(map[string]int64)
	if !ok {
		return ok
	}
	_, ok = secrets[key]
	if delAfterVerfiy && ok {
		delete(secrets, key)
		ctx.Session().Set(captchaGoSessionKey, secrets)
	}
	return ok
}
