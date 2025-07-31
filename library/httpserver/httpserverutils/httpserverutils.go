package httpserverutils

import "github.com/webx-top/echo"

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
