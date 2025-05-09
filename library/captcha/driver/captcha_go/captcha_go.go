package captcha_go

import (
	"html/template"

	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func newCaptchaGo() captchaLib.ICaptcha {
	return &captchaGo{}
}

var (
	cssURLs = []string{
		`/js/captchago/css/style.css`,
	}
	jsURLs = map[string]map[string][]string{
		`click`: {
			``: {
				`/js/captchago/js/common.js`,
				`/js/captchago/js/click.js`,
				`/js/captchago/js/jquery.captcha.js`,
			},
		},
		`rotate`: {
			``: {
				`/js/captchago/js/common.js`,
				`/js/captchago/js/rotate.js`,
				`/js/captchago/js/jquery.captcha.js`,
			},
		},
		`slide`: {
			`basic`: {
				`/js/captchago/js/common.js`,
				`/js/captchago/js/slide-basic.js`,
				`/js/captchago/js/jquery.captcha.js`,
			},
			`region`: {
				`/js/captchago/js/common.js`,
				`/js/captchago/js/slide-region.js`,
				`/js/captchago/js/jquery.captcha.js`,
			},
		},
	}
)

type captchaGo struct {
	driver    string
	cType     string
	captchaID string
	cfg       echo.H
}

func (c *captchaGo) Init(opt echo.H) error {
	c.driver = opt.String(`driver`, `click`)
	c.cfg = opt.GetStore(c.driver)
	c.cType = c.cfg.String(`type`, `basic`)
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *captchaGo) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	initedKey := `CaptchaGoJSInited`
	var loadResource bool
	if !ctx.Internal().Bool(initedKey) {
		ctx.Internal().Set(initedKey, true)
		loadResource = true
	}
	keysValues = append(keysValues, `loadResource`, loadResource)
	return c.render(ctx, templatePath, keysValues...)
}

func (c *captchaGo) render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues...)
	options.Set("driver", c.driver)
	options.Set("type", c.cType)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	options.Set("captchaID", c.captchaID)
	if !options.Has("captchaName") {
		options.Set("captchaName", "captchaGo")
	}
	options.Set("jsURLs", c.getJSURLs())
	options.Set("cssURLs", cssURLs)
	return captchaLib.RenderTemplate(ctx, captchaLib.TypeGo, templatePath, options)
}

func (c *captchaGo) getJSURLs() []string {
	g, y := jsURLs[c.driver]
	if !y {
		return []string{}
	}
	v, y := g[c.cType]
	if y {
		return v
	}
	return g[``]
}

func (c *captchaGo) Verify(ctx echo.Context, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	var idGet func(name string) []string
	if len(captchaIdent) > 0 {
		idGet = func(_ string) []string {
			return ctx.FormValues(captchaIdent[0])
		}
	} else {
		idGet = ctx.FormValues
	}
	id := idGet(captchaName)
	if len(id) == 0 {
		captchaName = "captchaGo"
		id = idGet(captchaName)
	}
	if len(id) == 0 { // 为空说明表单没有显示验证码输入框，此时返回验证码信息供前端显示
		data := captchaLib.GenCaptchaError(ctx, nil, captchaName, c.MakeData(ctx, hostAlias, captchaName))
		return data.SetInfo(ctx.T(`行为验证码显示失败`), captchaLib.ErrCaptchaIdMissing.Code.Int())
	}
	if len(id[0]) == 0 {
		return ctx.Data().SetInfo(ctx.T(`请进行行为验证`), captchaLib.ErrCaptchaCodeRequired.Code.Int()).SetZone(captchaName)
	}
	if !captchaGoVerifySuccessKey(ctx, id[0], true) {
		return ctx.Data().SetInfo(ctx.T(`行为验证未通过，请重试`), captchaLib.ErrCaptcha.Code.Int()).SetZone(captchaName)
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *captchaGo) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("driver", c.driver)
	data.Set("type", c.cType)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	data.Set("captchaType", captchaLib.TypeGo)
	data.Set("captchaID", c.captchaID)
	data.Set("jsURLs", c.getJSURLs())
	data.Set("cssURLs", cssURLs)
	jsInit := c.render(ctx, `partial_jsinit`, `captchaName`, name)
	data.Set("jsInit", jsInit)
	htmlCode := c.Render(ctx, `partial_main`, `captchaName`, name)
	data.Set("html", htmlCode)
	data.Set("captchaName", name)
	return data
}
