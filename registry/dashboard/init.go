package dashboard

import (
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/echo"
)

func init() {
	httpserver.Backend.Dashboard.Cards.Add(0, (&dashboard.Card{
		IconName:  `fa-user`,
		IconColor: `success`,
		Short:     `USERS`,
		Name:      `用户数量`,
		Summary:   ``,
	}).SetContentGenerator(func(ctx echo.Context) interface{} {
		//用户统计
		userMdl := model.NewUser(ctx)
		userCount, _ := userMdl.Count(nil)
		return userCount
	}))
	httpserver.Backend.Dashboard.TopButtons.Add(0,
		&dashboard.Button{
			Tmpl: `manager/topbutton/donation`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/language`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/source`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/bug-report`,
		},
	)
}
