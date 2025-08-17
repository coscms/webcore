package route

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
)

type IRegister interface {
	Echo() *echo.Echo
	Routes() []*echo.Route
	Logger() logger.Logger
	Prefix() string
	SetPrefix(prefix string) IRegister
	SetSkipper(f func() bool) IRegister
	Skipped() bool
	MakeHandler(handler interface{}, requests ...interface{}) echo.Handler
	MetaHandler(m echo.H, handler interface{}, requests ...interface{}) echo.Handler
	MetaHandlerWithRequest(m echo.H, handler interface{}, request interface{}, methods ...string) echo.Handler
	HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler
	AddGroupNamer(namers ...func(string) string) IRegister
	SetGroupNamer(namers ...func(string) string) IRegister
	SetRootGroup(groupName string) IRegister
	RootGroup() string
	Apply() IRegister
	Pre(middlewares ...interface{}) IRegister
	PreUse(middlewares ...interface{}) IRegister
	Use(middlewares ...interface{}) IRegister
	PreToGroup(groupName string, middlewares ...interface{}) IRegister
	PreUseToGroup(groupName string, middlewares ...interface{}) IRegister
	UseToGroup(groupName string, middlewares ...interface{}) IRegister
	Register(fn func(echo.RouteRegister)) IRegister
	RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) MetaSetter
	Host(hostName string, middlewares ...interface{}) Hoster
	Clear() IRegister
}
