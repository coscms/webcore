package model

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

const (
	AuthTypePassword = `password`
	AuthTypeOauth2   = `oauth2`
)

type SafeItemInfo struct {
	Step        uint
	ConfigTitle string
	ConfigRoute string
	hide        func() bool
}

func (s SafeItemInfo) IsZero() bool {
	return s.Step == 0
}

func (s *SafeItemInfo) SetHider(fn func() bool) *SafeItemInfo {
	s.hide = fn
	return s
}

func (s SafeItemInfo) IsHide() bool {
	if s.hide == nil {
		return false
	}
	return s.hide()
}

// var SafeItems = echo.NewKVData().
// 	Add(`gauth_bind`, `两步验证`).
// 	Add(`password`, `修改密码`)

var SafeItems = echo.NewKVData().
	Add(`google`, echo.T(`两步验证`), echo.KVOptX(SafeItemInfo{
		Step: 2, ConfigTitle: echo.T(`两步验证`), ConfigRoute: `gauth_bind`,
	})).
	Add(`password`, echo.T(`密码登录`), echo.KVOptX(SafeItemInfo{
		Step: 1, ConfigTitle: echo.T(`修改密码`), ConfigRoute: `password`,
	}))

var emptySafeItemInfo = SafeItemInfo{}

func RegisterSafeItem(itemType, itemTitle string, info SafeItemInfo, extra ...echo.H) {
	var _extra echo.H
	if len(extra) > 0 {
		_extra = extra[0]
	}
	SafeItems.Add(itemType, itemTitle, echo.KVOptX(info), echo.KVOptH(_extra))
}

func GetSafeItem(itemType string) SafeItemInfo {
	item := SafeItems.GetItem(itemType)
	if item == nil {
		return emptySafeItemInfo
	}
	v, _ := item.X.(SafeItemInfo)
	return v
}

func ListSafeItemsByStep(step uint, exclude ...string) []echo.KV {
	items := SafeItems.Slice()
	result := make([]echo.KV, 0, len(items))
	for _, item := range items {
		v, _ := item.X.(SafeItemInfo)
		if v.Step == step && !com.InSlice(item.K, exclude) {
			result = append(result, *item)
		}
	}
	return result
}
