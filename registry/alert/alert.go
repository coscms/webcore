package alert

import (
	"github.com/coscms/webcore/library/imbot"
	_ "github.com/coscms/webcore/library/imbot/dingding"
	_ "github.com/coscms/webcore/library/imbot/workwx"
	"github.com/webx-top/echo"
)

var (
	// RecipientTypes 收信类型
	RecipientTypes = echo.NewKVData()

	// RecipientPlatforms 收信平台
	RecipientPlatforms             = echo.NewKVData()
	RecipientPlatformWebhookCustom = `custom`

	// Topics 告警专题
	Topics = echo.NewKVData()
)

func init() {
	RecipientTypes.Add(`email`, `email`)
	RecipientTypes.Add(`webhook`, `webhook`)
	for name, mess := range imbot.Messagers() {
		RecipientPlatforms.Add(name, mess.Label)
	}
	RecipientPlatforms.Add(RecipientPlatformWebhookCustom, echo.T(`自定义`))

	//Topics.Add(`test`, `测试`)
}
