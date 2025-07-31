package httpserver

import (
	"strings"

	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/route"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webcore/library/httpserver/httpserverutils"
)

const (
	KindBackend   = httpserverutils.KindBackend
	KindFrontend  = httpserverutils.KindFrontend
	ServerKindKey = httpserverutils.ServerKindKey
)

var (
	Backend = New(KindBackend).
		SetRouter(route.NewRegister(defaults.Default).AddGroupNamer(groupNamer)).
		SetNavigate(navigate.NewProjectNavigates(KindBackend, `nging`))
	Frontend = New(KindFrontend).SetNavigate(navigate.NewProjectNavigates(KindFrontend, `webx`))
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
