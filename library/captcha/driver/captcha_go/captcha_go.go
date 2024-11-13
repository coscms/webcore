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
	initialized.Do(initialize)
	return &captchaGo{jsURL: `1`}
}

type captchaGo struct {
	driver    string
	cType     string
	jsURL     string
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
	options := tplfunc.MakeMap(keysValues)
	options.Set("driver", c.driver)
	options.Set("type", c.cType)
	initedKey := `CaptchaGoJSInited.` + c.driver
	var jsURL string
	if !ctx.Internal().Bool(initedKey) {
		ctx.Internal().Set(initedKey, true)
		jsURL = c.jsURL
	}
	options.Set("jsURL", jsURL)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	options.Set("captchaID", c.captchaID)
	if !options.Has("captchaName") {
		options.Set("captchaName", "captchaGo")
	}
	return captchaLib.RenderTemplate(ctx, captchaLib.TypeGo, templatePath, options)
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
		return ctx.Data().SetInfo(ctx.T(`行为验证码显示失败`), captchaLib.ErrCaptchaIdMissing.Code.Int())
	}
	if len(id[0]) == 0 {
		return ctx.Data().SetInfo(ctx.T(`请进行行为验证`), captchaLib.ErrCaptcha.Code.Int())
	}
	if !captchaGoVerifySuccessKey(ctx, id[0], true) {
		data := captchaLib.GenCaptchaError(ctx, nil, captchaName, c.MakeData(ctx, hostAlias, captchaName))
		return data.SetInfo(ctx.T(`行为验证未通过，请重试`), captchaLib.ErrCaptcha.Code.Int())
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
	htmlCode := c.Render(ctx, ``)
	data.Set("html", htmlCode)
	data.Set("captchaName", name)
	return data
}
