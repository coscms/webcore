package oauth2client

import (
	"github.com/coscms/webcore/model"
)

func init() {
	item := model.SafeItemInfo{
		Step: 1, ConfigTitle: `账号绑定`, ConfigRoute: `oauth`,
	}
	item.SetHider(func() bool {
		accounts := GetOAuthAccounts()
		return len(accounts) == 0
	})
	model.RegisterSafeItem(`oauth`, `oAuth登录`, item)
}
