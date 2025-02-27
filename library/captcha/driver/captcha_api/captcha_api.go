package captcha_api

import (
	"html/template"
	"strconv"

	"github.com/admpub/captcha-go"
	"github.com/admpub/log"
	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/com/formatter"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func newCaptchaAPI() captchaLib.ICaptcha {
	return &captchaAPI{}
}

type captchaAPI struct {
	provider  string
	endpoint  captcha.Endpoint
	siteKey   string
	verifier  *captcha.SimpleCaptchaVerifier
	jsURL     string
	captchaID string
	cfg       echo.H
}

func (c *captchaAPI) Init(opt echo.H) error {
	c.provider = opt.String(`provider`)
	c.cfg = opt.GetStore(c.provider)
	c.siteKey = c.cfg.String(`siteKey`)
	switch c.provider {
	case `turnstile`:
		c.endpoint = captcha.CloudflareTurnstile
		c.jsURL = `https://challenges.cloudflare.com/turnstile/v0/api.js`
	default:
		c.endpoint = captcha.GoogleRecaptcha
		c.jsURL = `https://www.recaptcha.net/recaptcha/api.js?render=` + c.siteKey
	}
	captchaSecret := c.cfg.String(`secret`)
	//echo.Dump(echo.H{`siteKey`: c.siteKey, `secret`: captchaSecret})
	v := captcha.NewCaptchaVerifier(c.endpoint, captchaSecret)
	c.verifier = &captcha.SimpleCaptchaVerifier{
		Verifier: *v,
	}
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *captchaAPI) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("siteKey", c.siteKey)
	options.Set("endpoint", c.endpoint)
	options.Set("provider", c.provider)
	var jsURL string
	c.captchaID = com.RandomAlphanumeric(16)
	initedKey := `CaptchaJSInited.` + c.provider
	if !ctx.Internal().Bool(initedKey) {
		ctx.Internal().Set(initedKey, true)
		jsURL = c.jsURL
	}
	options.Set("jsURL", jsURL)
	options.Set("captchaID", c.captchaID)
	if !options.Has("captchaName") {
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			options.Set("captchaName", "cf-turnstile-response")
		default:
			options.Set("captchaName", "g-recaptcha-response")
		}
	}
	if len(templatePath) == 0 {
		var htmlContent string
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `"></script>`
			}
			locationID := `turnstile-` + c.captchaID
			htmlContent += `<input type="hidden" name="captchaId" value="` + c.captchaID + `" />`
			var theme string
			if ctx.Cookie().Get(`ThemeColor`) == `dark` {
				theme = `dark`
			} else {
				theme = `light`
			}
			htmlContent += `<div class="cf-turnstile" id="` + locationID + `" data-sitekey="` + c.siteKey + `" data-theme="` + theme + `"></div>`
			htmlContent += `<input type="hidden" id="` + locationID + `-extend" disabled />`
			htmlContent += `<script>
window.addEventListener('load', function(){
    var id='#` + locationID + `', $box=$(id);
	$box.closest('.input-group-addon').addClass('xxs-padding-top').prev('input').remove();
	var $form=$box.closest('form');
	$form.on('submit',function(e){
		if($box.data('lastGeneratedAt')>(new Date()).getTime()-290) {
			$box.data('lastGeneratedAt',0);
			return true;
		}
		window.setTimeout(function(){turnstile.reset(id);},1000);
		$box.data('lastGeneratedAt',(new Date()).getTime());
	});
    if(!$('body').data('tarnstileInited')){$('body').data('tarnstileInited',true);return;}
    if($box.children('div').length>0)return;
    turnstile.render(id);
})
</script>`
		default:
			locationID := `recaptcha-` + c.captchaID
			htmlContent = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" />`
			htmlContent += `<input type="hidden" id="` + locationID + `" name="g-recaptcha-response" value="" />`
			htmlContent += `<input type="hidden" id="` + locationID + `-extend" disabled />`
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `"></script>`
			}
			htmlContent += `<script>
window.addEventListener('load', function(){
	var id='#` + locationID + `';
	grecaptcha.ready(function() {
	  grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		$(id).val(token);
		$(id).data('lastGeneratedAt',(new Date()).getTime());
	  });
	});
	var igrp=$(id).closest('.input-group');
	if(igrp.length>0){
		igrp.hide();
		if(igrp.parent().hasClass('form-group')) igrp.parent().hide();
	}
	$(id).closest('.input-group-addon').prev('input').remove();
	var $form=$(id).closest('form');
	var $submit=$form.find(':submit');
	$submit.on('click',function(e){
		if($(id).val() && $(id).data('lastGeneratedAt')>(new Date()).getTime()-110) {
			$(id).data('lastGeneratedAt',0);
			return true;
		}
		var $this=$(this);
		e.preventDefault();
		grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		  $(id).val(token);
		  $(id).data('lastGeneratedAt',(new Date()).getTime());
		  $this.trigger('click');
		});
	});
})
</script>`
		}
		return template.HTML(htmlContent)
	}
	return captchaLib.RenderTemplate(ctx, captchaLib.TypeAPI, templatePath, options)
}

func (c *captchaAPI) Verify(ctx echo.Context, hostAlias string, captchaName string, _ ...string) echo.Data {
	var name string
	switch c.endpoint {
	case captcha.CloudflareTurnstile:
		name = `cf-turnstile-response`
	default:
		name = `g-recaptcha-response`
		if c.cfg.Has(`minScore`) {
			c.verifier.MinScore = c.cfg.Float32(`minScore`)
		} else {
			c.verifier.MinScore = 0.5
		}
		c.verifier.ExpectedAction = ctx.Form(`captchaAction`, `submit`)
	}
	vcode := ctx.FormValues(captchaName)
	if len(vcode) == 0 {
		captchaName = name
		vcode = ctx.FormValues(captchaName)
	}
	c.captchaID = ctx.Formx(`captchaId`).String()
	if len(c.captchaID) == 0 {
		return captchaLib.GenCaptchaError(ctx, captchaLib.ErrCaptchaIdMissing, captchaName, c.MakeData(ctx, hostAlias, captchaName))
	}
	if len(vcode) == 0 || len(vcode[0]) == 0 { // 为空说明没有验证码
		data := captchaLib.GenCaptchaError(ctx, nil, captchaName, c.MakeData(ctx, hostAlias, captchaName))
		return data.SetInfo(ctx.T(`请先进行人机验证`), captchaLib.ErrCaptchaCodeRequired.Code.Int())
	}
	token := vcode[0]
	c.verifier.ExpectedHostname = ctx.Domain()
	var clientIP string
	if c.cfg.Bool(`verifyIP`) {
		clientIP = ctx.RealIP()
	}
	resp, ok, err := c.verifier.VerifyActionWithResponse(token, clientIP, c.verifier.ExpectedAction)
	if err != nil {
		return captchaLib.GenCaptchaError(ctx, err, captchaName, c.MakeData(ctx, hostAlias, captchaName))
	}
	if !ok {
		log.Warnf(`failed to captchaAPI.Verify: %s`, formatter.AsStringer(resp))
		data := captchaLib.GenCaptchaError(ctx, nil, captchaName, c.MakeData(ctx, hostAlias, captchaName))
		return data.SetInfo(ctx.T(`未能通过人机验证，请重试`), captchaLib.ErrCaptcha.Code.Int())
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *captchaAPI) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("siteKey", c.siteKey)
	data.Set("provider", c.provider)
	data.Set("jsURL", c.jsURL)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	data.Set("captchaType", captchaLib.TypeAPI)
	data.Set("captchaID", c.captchaID)
	var jsInit, jsCallback string
	var locationID string
	var htmlCode string
	var captchaName string
	switch c.provider {
	case `turnstile`:
		captchaName = `cf-turnstile-response`
		locationID = `turnstile-` + c.captchaID
		jsInit = `(function(){ typeof(turnstile)!='undefined' && turnstile.render('#` + locationID + `'); })();`
		jsCallback = `function(callback){
	callback && callback();
	window.setTimeout(function(){turnstile.reset('#` + locationID + `');},1000);
}`
		var theme string
		if ctx.Cookie().Get(`ThemeColor`) == `dark` {
			theme = `dark`
		} else {
			theme = `light`
		}
		htmlCode = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" /><div class="cf-turnstile text-center" id="turnstile-` + c.captchaID + `" data-sitekey="` + c.siteKey + `" data-theme="` + theme + `"></div>`
	default:
		captchaName = `g-recaptcha-response`
		locationID = `recaptcha-` + c.captchaID
		defaultTips := strconv.Quote(ctx.T(`加载成功，请点击“提交”按钮继续`))
		jsInit = `(function(){
var f=function(){
	if(typeof(grecaptcha)=='undefined'){setTimeout(f,200);return;}
	grecaptcha.ready(function() {
		var id='#` + locationID + `';
		grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
			$(id).val(token);
			$(id).data('lastGeneratedAt',(new Date()).getTime());
			var $loading=$('#` + locationID + `-loading');
			if($loading.length>0){
				var t=$loading.data('success-tips')||` + defaultTips + `;
				$loading.html('<i class="fa fa-check text-success"></i> '+t);
			}
		});
	});
};
f();
})();`
		jsCallback = `function(callback){
	grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		var id='#` + locationID + `';
		$(id).val(token);
		$(id).data('lastGeneratedAt',(new Date()).getTime());
		callback && callback(token);
	});
}`
		htmlCode = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" /><input type="hidden" id="recaptcha-` + c.captchaID + `" name="g-recaptcha-response" value="" />`
	}
	data.Set("jsCallback", jsCallback)
	data.Set("jsInit", jsInit)
	data.Set("locationID", locationID)
	data.Set("html", htmlCode)
	data.Set("captchaIdent", `captchaId`)
	data.Set("captchaName", captchaName)
	return data
}
