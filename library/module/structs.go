package module

import (
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/route"
)

func NewNavigate() *Navigate {
	return &Navigate{
		Backend:  httpserver.Backend.Navigate,
		Frontend: httpserver.Frontend.Navigate,
	}
}

type Navigate struct {
	Backend  *navigate.ProjectNavigates
	Frontend *navigate.ProjectNavigates
}

func NewDashboard() *Dashboard {
	return &Dashboard{
		Backend:  httpserver.Backend.Dashboard,
		Frontend: httpserver.Frontend.Dashboard,
	}
}

type Dashboard struct {
	Backend  *dashboard.Dashboard
	Frontend *dashboard.Dashboard
}

func NewRouter() *Router {
	return &Router{
		Backend:  httpserver.Backend.Router,
		Frontend: httpserver.Frontend.Router,
	}
}

type Router struct {
	Backend  route.IRegister
	Frontend route.IRegister
}
