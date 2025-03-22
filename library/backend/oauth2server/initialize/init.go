package initialize

import (
	"encoding/gob"

	"github.com/admpub/oauth2/v4"
	"github.com/coscms/oauth2s"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend/oauth2server"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/registry/route"
)

func Init(loginHandler func(echo.Context) error, logoutHandler func(echo.Context) error) {
	func() {
		defer recover()
		gob.Register(map[string][]string{})
	}()
	oauth2server.RoutePrefix = `/oauth2`
	route.Register(func(r echo.RouteRegister) {
		oauth2server.Debug = !config.FromFile().Sys.IsEnv(`prod`)
		var tokenStore oauth2.TokenStore
		oauth2server.Default.Init(
			oauth2s.JWTMethod(nil),
			//oauth2s.JWTKey([]byte(config.FromFile().Cookie.HashKey)),
			//oauth2s.JWTMethod(jwt.SigningMethodHS512),
			oauth2s.ClientStore(oauth2server.DefaultAppClient),
			oauth2s.SetStore(tokenStore),
			oauth2s.SetHandler(&oauth2s.HandlerInfo{
				PasswordAuthorization: oauth2server.PasswordAuthorizationHandler,
				UserAuthorize:         oauth2server.UserAuthorizeHandler,
				InternalError:         oauth2server.InternalErrorHandler,
				ResponseError:         oauth2server.ResponseErrorHandler,
				RefreshingScope:       oauth2server.RefreshingScopeHandler,
				RefreshingValidation:  oauth2server.RefreshingValidationHandler,
			}),
		)
		oauth2server.Route(r, loginHandler, logoutHandler)
	})
}
