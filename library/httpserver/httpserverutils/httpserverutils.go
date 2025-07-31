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
