package navigate

import (
	"slices"

	"github.com/webx-top/echo"
)

type Group struct {
	Icon  string
	Label string
	Group string
}

type Groups struct {
	g []Group
	k map[string]int
}

func (g *Groups) Slice() []Group {
	return g.g
}

var navGroups = map[string]*Groups{}

func init() {
	RegisterGroup(`backend.top.manager`, []Group{
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
	}...)
}

func RegisterGroup(name string, groups ...Group) {
	_, y := navGroups[name]
	if !y {
		navGroups[name] = &Groups{
			k: map[string]int{},
		}
	}
	for _, group := range groups {
		i, y := navGroups[name].k[group.Group]
		if y {
			navGroups[name].g[i] = group
			continue
		}
		navGroups[name].k[group.Group] = len(navGroups[name].g)
		navGroups[name].g = append(navGroups[name].g, group)
	}
}

func HasGroup(name string, group string) bool {
	_, y := navGroups[name]
	if !y {
		return false
	}
	_, y = navGroups[name].k[group]
	return y
}

func UnregisterGroup(name string, group string) {
	v, y := navGroups[name]
	if !y {
		return
	}
	navGroups[name].g = slices.DeleteFunc(v.g, func(e Group) bool {
		if e.Group == group {
			delete(navGroups[name].k, group)
			return true
		}
		return false
	})
}

func GetGroups(name string) []Group {
	v, y := navGroups[name]
	if !y || v == nil {
		return nil
	}
	return v.g
}

func GetGroupIndexes(name string) map[string]int {
	v, y := navGroups[name]
	if !y || v == nil {
		return nil
	}
	return v.k
}
