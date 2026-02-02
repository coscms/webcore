package pipe

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	modelFile "github.com/coscms/webcore/model/file"
	"github.com/coscms/webcore/registry/upload/driver"
	uploadClient "github.com/webx-top/client/upload"
)

func init() {
	Register(`_queryMd5`, queryMd5) // 以下划线开始表示这个独立的功能
}

// queryMd5 查询MD5
func queryMd5(ctx echo.Context, _ driver.Storer, _ uploadClient.Results, data map[string]interface{}) error {
	md5 := ctx.Form(`md5`)
	if len(md5) == 0 {
		return ctx.NewError(code.InvalidParameter, `MD5值不正确`).SetZone(`md5`)
	}
	m := modelFile.NewFile(ctx)
	err := m.GetByMd5(md5)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}

	data[`file`] = m.ViewUrl
	data[`name`] = m.Name
	return nil
}
