package captcha_go

import "github.com/coscms/webcore/library/captcha"

func init() {
	captcha.Register(captcha.TypeGo, newCaptchaGo)
}
