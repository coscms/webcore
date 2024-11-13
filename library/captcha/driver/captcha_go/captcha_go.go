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
	return &captchaGo{jsURL: `1`}
}

type captchaGo struct {
	provider  string
	jsURL     string
	captchaID string
	cfg       echo.H
}

func (c *captchaGo) Init(opt echo.H) error {
	c.provider = opt.String(`provider`)
	c.cfg = opt.GetStore(c.provider)
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *captchaGo) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("provider", c.provider)
	initedKey := `CaptchaGoJSInited.` + c.provider
	var jsURL string
	if !ctx.Internal().Bool(initedKey) {
		ctx.Internal().Set(initedKey, true)
		jsURL = c.jsURL
	}
	options.Set("jsURL", jsURL)
	c.captchaID = com.RandomAlphanumeric(16)
	options.Set("captchaID", c.captchaID)
	return captchaLib.RenderTemplate(ctx, captchaLib.TypeGo, templatePath, options)
}

func (c *captchaGo) Verify(ctx echo.Context, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	var idGet func(name string, defaults ...string) string
	if len(captchaIdent) > 0 {
		idGet = func(_ string, defaults ...string) string {
			return ctx.Form(captchaIdent[0], defaults...)
		}
	} else {
		idGet = ctx.Form
	}
	id := idGet("captchaGo")
	if len(id) == 0 { // 为空说明表单没有显示验证码输入框，此时返回验证码信息供前端显示
		return ctx.Data().SetError(captchaLib.ErrCaptchaIdMissing)
	}
	if !captchaGoVerifySuccessKey(ctx, id, true) {
		return captchaLib.GenCaptchaError(ctx, captchaLib.ErrCaptcha, captchaName, c.MakeData(ctx, hostAlias, captchaName))
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *captchaGo) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("provider", c.provider)
	data.Set("jsURL", c.jsURL)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	data.Set("captchaType", captchaLib.TypeGo)
	data.Set("captchaID", c.captchaID)
	locationID := `captchago-` + c.captchaID
	var jsCallback string
	jsInit := `(function(){
    var baseURL=IS_BACKEND?BACKEND_URL:FRONTEND_URL;
    $('#` + locationID + `').captcha({
        api:baseURL+'/captchago',
        success:function(){
        },
        error:function(){
        }
    });
	})()`
	htmlCode := ``
	var captchaName string
	data.Set("jsCallback", jsCallback)
	data.Set("jsInit", jsInit)
	data.Set("locationID", locationID)
	data.Set("html", htmlCode)
	data.Set("captchaIdent", `captchaGo`)
	data.Set("captchaName", captchaName)
	return data
}
