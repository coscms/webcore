package oauth2server

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/middleware"
	"github.com/webx-top/echo"
)

var (
	// RoutePrefix 路由前缀
	RoutePrefix string
)

func Route(g echo.RouteRegister, loginHandler func(echo.Context) error, logoutHandler func(echo.Context) error) {
	if len(RoutePrefix) > 0 {
		g = g.Group(RoutePrefix)
	}
	g.Route(`GET,POST`, `/authorize`, authorizeHandler, middleware.AuthCheck).SetMetaKV(httpserver.PermGuestKV())              // (step.1)(step.4)
	g.Route(`GET,POST`, `/login`, makeLoginHandler(loginHandler), middleware.AuthCheck).SetMetaKV(httpserver.PermGuestKV())    // (step.2) 登录页面
	g.Route(`GET,POST`, `/auth`, authHandler, middleware.AuthCheck).SetMetaKV(httpserver.PermPublicKV())                       // (step.3) 授权页面
	g.Route(`GET,POST`, `/logout`, makeLogoutHandler(logoutHandler), middleware.AuthCheck).SetMetaKV(httpserver.PermGuestKV()) // 退出登录
	g.Route(`GET,POST`, `/token`, tokenHandler, middleware.AuthCheck).SetMetaKV(httpserver.PermGuestKV())
	g.Route(`GET,POST`, `/profile`, profileHandler, middleware.AuthCheck).SetMetaKV(httpserver.PermGuestKV()) // 获取用户个人资料
}
