package httpserverutils

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session"
)

const (
	KindBackend   = `backend`
	KindFrontend  = `frontend`
	ServerKindKey = `HTTP_SERVER_KIND`
)

func GetServerKind(e *echo.Echo) string {
	return e.Extra().String(ServerKindKey)
}

func GetServerKindByContext(ctx echo.Context) string {
	return GetServerKind(ctx.Echo())
}

func IsFrontend(e *echo.Echo) bool {
	return GetServerKind(e) == KindFrontend
}

func IsBackend(e *echo.Echo) bool {
	return GetServerKind(e) == KindBackend
}

func IsFrontendContext(ctx echo.Context) bool {
	return GetServerKindByContext(ctx) == KindFrontend
}

func IsBackendContext(ctx echo.Context) bool {
	return GetServerKindByContext(ctx) == KindBackend
}

// CookieMaxAge 允许设置的Cookie最大有效时长(单位:秒)
var CookieMaxAge = 86400 * 365

func RememberSession(c echo.Context) int {
	remember := c.Form(`remember`)
	var maxAge int
	if len(remember) > 0 {
		if remember == `forever` {
			maxAge = CookieMaxAge
		} else {
			duration, err := com.ParseTimeDuration(remember)
			if err == nil {
				maxAge = int(duration.Seconds())
			}
		}
		if maxAge > CookieMaxAge {
			maxAge = CookieMaxAge
		}
		session.RememberMaxAge(c, maxAge)
	}
	return maxAge
}
