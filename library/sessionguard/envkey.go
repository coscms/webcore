package sessionguard

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func EnvKey(ctx echo.Context, cfg SessionGuardConfig) string {
	if cfg.IgnoreBrowserIP && cfg.IgnoreBrowserUA {
		return `T` + strconv.FormatInt(time.Now().Unix(), 10)
	}
	if cfg.IgnoreBrowserUA {
		return com.Md5(ctx.RealIP()) + `T` + strconv.FormatInt(time.Now().Unix(), 10)
	}
	if cfg.IgnoreBrowserIP {
		return com.Md5(`|`+ctx.Request().UserAgent()) + `T` + strconv.FormatInt(time.Now().Unix(), 10)
	}
	return com.Md5(ctx.RealIP()+`|`+ctx.Request().UserAgent()) + `T` + strconv.FormatInt(time.Now().Unix(), 10)
}

func (cfg *SessionGuardConfig) VerifyEnvKey(ctx echo.Context, ownerType string, ownerIdOrName interface{}, envKey string, dontCheckTime bool) error {
	pos := strings.LastIndex(envKey, `T`)
	end := len(envKey) - 1
	if pos < 1 || pos >= end {
		return ctx.NewError(code.InvalidParameter, `凭证无效`).SetZone(`envKey`)
	}
	if !dontCheckTime {
		ts, err := strconv.ParseInt(envKey[pos+1:], 10, 64)
		if err != nil {
			return ctx.NewError(code.DataFormatIncorrect, `凭证解析失败`).SetZone(`envKey`)
		}
		if time.Now().Unix()-ts > 300 {
			return ctx.NewError(code.DataHasExpired, `凭证已经过期，请刷新页面后重新操作`).SetZone(`envKey`)
		}
	}
	if cfg.IgnoreBrowserUA && cfg.IgnoreBrowserIP {
		return nil
	}
	if cfg.IgnoreBrowserUA {
		currentIP := ctx.RealIP()
		lastIPMd5 := envKey[0:pos]
		currentIPMd5 := com.Md5(currentIP)
		if currentIPMd5 != lastIPMd5 {
			log.Debugf(`[%s:%v]ip mismatched: %q(last) != %q(now)`, ownerType, ownerIdOrName, lastIPMd5, currentIPMd5)
			return ctx.NewError(code.DataStatusIncorrect, `凭证来源不符合要求`).SetZone(`envKey`)
		}
	} else if cfg.IgnoreBrowserIP {
		currentUA := ctx.Request().UserAgent()
		lastUAMd5 := envKey[0:pos]
		currentUAMd5 := com.Md5(`|` + currentUA)
		if currentUAMd5 != lastUAMd5 {
			log.Debugf(`[%s:%v]userAgent mismatched: %q(last) != %q(now)`, ownerType, ownerIdOrName, lastUAMd5, currentUAMd5)
			return ctx.NewError(code.DataStatusIncorrect, `凭证来源不符合要求`).SetZone(`envKey`)
		}
	} else {
		currentIP := ctx.RealIP()
		lastIPAndUAMd5 := envKey[0:pos]
		if com.Md5(currentIP+`|`+ctx.Request().UserAgent()) != lastIPAndUAMd5 {
			return ctx.NewError(code.DataStatusIncorrect, `凭证来源不符合要求`).SetZone(`envKey`)
		}
	}
	return nil
}

type PE struct {
	Password  string `json:"p"`
	EnvKey    string `json:"e"`
	Timestamp int64  `json:"t"`
}

func (p *PE) Verify(ctx echo.Context, ownerType string, ownerIdOrName interface{}) error {
	cfg := GetConfig()
	sessionGuardCfg := cfg.SessionGuardConfig
	if ownerType != `user` {
		if cfg.Frontend == nil {
			return nil
		}
		sessionGuardCfg = *cfg.Frontend
	}
	encryptedPasswordExpires := sessionGuardCfg.Expires
	if encryptedPasswordExpires <= 0 {
		encryptedPasswordExpires = 300
	}
	if time.Now().Unix()-p.Timestamp > encryptedPasswordExpires {
		return ctx.NewError(code.DataHasExpired, `凭证已经失效`).SetZone(`envKey`)
	}
	return sessionGuardCfg.VerifyEnvKey(ctx, ownerType, ownerIdOrName, p.EnvKey, true)
}

func Unpack(ctx echo.Context, ownerType string, ownerIdOrName interface{}, encrypted string) (password string, err error) {
	if !strings.HasPrefix(encrypted, `{`) {
		return encrypted, nil
	}
	pe := &PE{}
	jerr := json.Unmarshal(com.Str2bytes(encrypted), pe)
	if jerr != nil {
		err = ctx.NewError(code.DataFormatIncorrect, `密码拆包失败`)
		return
	}
	return pe.Password, pe.Verify(ctx, ownerType, ownerIdOrName)
}
