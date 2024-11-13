package captcha_api

import "github.com/coscms/webcore/library/captcha"

func init() {
	captcha.Register(captcha.TypeAPI, newCaptchaAPI)
}
