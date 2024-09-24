package module

import (
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/route"
)

func NewNavigate() Navigate {
	return &paramNavigate{
		backend:  httpserver.Backend.Navigate,
		frontend: httpserver.Frontend.Navigate,
	}
}

type Navigate interface {
	Backend() *navigate.ProjectNavigates
	Frontend() *navigate.ProjectNavigates
}

type paramNavigate struct {
	backend  *navigate.ProjectNavigates
	frontend *navigate.ProjectNavigates
}

func (a *paramNavigate) Backend() *navigate.ProjectNavigates {
	return a.backend
}

func (a *paramNavigate) Frontend() *navigate.ProjectNavigates {
	return a.frontend
}

func NewDashboard() Dashboard {
	return &paramDashboard{
		backend:  httpserver.Backend.Dashboard,
		frontend: httpserver.Frontend.Dashboard,
	}
}

type Dashboard interface {
	Backend() *dashboard.Dashboard
	Frontend() *dashboard.Dashboard
}

type paramDashboard struct {
	backend  *dashboard.Dashboard
	frontend *dashboard.Dashboard
}

func (a *paramDashboard) Backend() *dashboard.Dashboard {
	return a.backend
}

func (a *paramDashboard) Frontend() *dashboard.Dashboard {
	return a.frontend
}

func NewRouter() Router {
	return &paramRouter{
		backend:  httpserver.Backend.Router,
		frontend: httpserver.Frontend.Router,
	}
}

type Router interface {
	Backend() route.IRegister
	Frontend() route.IRegister
}

type paramRouter struct {
	backend  route.IRegister
	frontend route.IRegister
}

func (a *paramRouter) Backend() route.IRegister {
	return a.backend
}

func (a *paramRouter) Frontend() route.IRegister {
	return a.frontend
}
