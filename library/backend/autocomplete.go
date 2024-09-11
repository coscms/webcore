package backend

import (
	"github.com/coscms/webcore/library/filemanager"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var onAutoCompletePath = []func(echo.Context) (bool, error){}

func OnAutoCompletePath(fn func(echo.Context) (bool, error)) {
	onAutoCompletePath = append(onAutoCompletePath, fn)
}

func FireAutoCompletePath(c echo.Context) (bool, error) {
	for _, fn := range onAutoCompletePath {
		ok, err := fn(c)
		if ok || err != nil {
			return true, err
		}
	}
	return false, nil
}

func AutoCompletePath(ctx echo.Context) error {
	user := User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	if ok, err := FireAutoCompletePath(ctx); ok || err != nil {
		return err
	}
	data := ctx.Data()
	prefix := ctx.Form(`query`)
	size := ctx.Formx(`size`, `10`).Int()
	var paths []string
	switch ctx.Form(`type`) {
	case `dir`:
		paths = filemanager.SearchDir(prefix, size)
	case `file`:
		paths = filemanager.SearchFile(prefix, size)
	default:
		paths = filemanager.Search(prefix, size)
	}
	data.SetData(paths)
	return ctx.JSON(data)
}
