package navigate

import (
	"slices"

	"github.com/webx-top/echo"
)

type Group struct {
	Icon  string
	Label string
	Group string
	Nav   *List
}

type GroupWithList struct {
	Icon  string
	Label string
	Group string
}

var navGroups = map[string][]Group{
	`backend.top.manager`: {
		{
			Icon:  `gear`,
			Label: echo.T(`基本设置`),
		},
		{
			Icon:  `gear`,
			Label: echo.T(`管理员账号`),
			Group: `admin`,
		},
		{
			Icon:  `file`,
			Label: echo.T(`附件管理`),
			Group: `file`,
		},
		{
			Icon:  `flask`,
			Label: echo.T(`缓存管理`),
			Group: `cache`,
		},
	},
}

func RegisterGroup(name string, groups ...Group) {
	v, y := navGroups[name]
	if !y {
		navGroups[name] = groups
	} else {
		v = append(v, groups...)
		navGroups[name] = v
	}
}

func UnregisterGroup(name string, group string) {
	v, y := navGroups[name]
	if y {
		navGroups[name] = slices.DeleteFunc(v, func(e Group) bool {
			return e.Group == group
		})
	}
}
