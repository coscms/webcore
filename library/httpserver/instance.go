package httpserver

import (
	"strings"

	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/route"
	"github.com/webx-top/echo/defaults"
)

const (
	KindBackend  = `backend`
	KindFrontend = `frontend`
)

var (
	Backend = New(KindBackend).
		SetRouter(route.NewRegister(defaults.Default).AddGroupNamer(groupNamer)).
		SetNavigate(navigate.NewProjectNavigates(`nging`))
	Frontend = New(KindFrontend).SetNavigate(navigate.NewProjectNavigates(`webx`))
	Servers  = &HTTPServers{
		Backend:  Backend,
		Frontend: Frontend,
	}
)

func groupNamer(group string) string {
	if len(group) == 0 {
		return group
	}
	if group == `@` {
		return ``
	}
	if strings.HasPrefix(group, `@`) {
		return group[1:]
	}
	return group
}

func Clear() {
	Servers.Clear()
}
