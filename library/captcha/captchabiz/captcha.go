package captchabiz

import (
	"html/template"

	"github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	_ "github.com/coscms/webcore/library/captcha/driver/captcha_api"
	_ "github.com/coscms/webcore/library/captcha/driver/captcha_go"
)

func GetCaptchaType(types ...string) (echo.H, string) {
	cfg := config.FromDB().GetStore(`captcha`)
	typ := cfg.String(`type`, captcha.TypeDefault)
	if len(types) > 0 && len(types[0]) > 0 {
		typ = types[0]
	}
	if len(typ) == 0 {
		typ = captcha.TypeDefault
	}
	return cfg, typ
}

func GetCaptchaEngine(ctx echo.Context, types ...string) (captcha.ICaptcha, error) {
	cfg, typ := GetCaptchaType(types...)
	create := captcha.Get(typ)
	if create == nil {
		if typ != captcha.TypeDefault {
			create = captcha.Get(captcha.TypeDefault)
			if create == nil {
				return nil, ctx.NewError(code.Unsupported, `不支持验证码类型: %s`, typ)
			}
		} else {
			return nil, ctx.NewError(code.Unsupported, `不支持验证码类型: %s`, typ)
		}
	}
	cpt := create()
	tcfg := cfg.Children(typ)
	err := cpt.Init(tcfg)
	return cpt, err
}

// VerifyCaptcha 验证码验证
func VerifyCaptcha(ctx echo.Context, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	cpt, err := GetCaptchaEngine(ctx)
	if err != nil {
		return ctx.Data().SetError(err)
	}
	return cpt.Verify(ctx, hostAlias, captchaName, captchaIdent...)
}

// VerifyCaptchaWithType 验证指定类型验证码
func VerifyCaptchaWithType(ctx echo.Context, captchaType string, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	cpt, err := GetCaptchaEngine(ctx, captchaType)
	if err != nil {
		return ctx.Data().SetError(err)
	}
	return cpt.Verify(ctx, hostAlias, captchaName, captchaIdent...)
}

// CaptchaForm 验证码输入表单
func CaptchaForm(ctx echo.Context, tmpl string, args ...interface{}) template.HTML {
	cpt, err := GetCaptchaEngine(ctx)
	if err != nil {
		return template.HTML(err.Error())
	}
	return cpt.Render(ctx, tmpl, args...)
}

// CaptchaFormWithType 生成指定类型的验证码输入表单
func CaptchaFormWithType(ctx echo.Context, captchaType string, tmpl string, args ...interface{}) template.HTML {
	cpt, err := GetCaptchaEngine(ctx, captchaType)
	if err != nil {
		return template.HTML(err.Error())
	}
	return cpt.Render(ctx, tmpl, args...)
}

func CheckEnable(typ string) func(h echo.Handler) echo.HandlerFunc {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := config.FromDB().GetStore(`captcha`)
			if cfg.String(`type`, captcha.TypeDefault) != typ {
				return echo.ErrNotFound
			}

			return h.Handle(c)
		}
	}
}
