package logcategory

import (
	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

type Other []string

var Categories = echo.NewKVxData[Other, any]().
	Add(log.DefaultLog.Category, echo.T(`Nging日志`), echo.KVxOptX[Other, any](Other{`websocket`, `watcher`})).
	Add(`db`, echo.T(`SQL日志`)).
	Add(`echo`, echo.T(`Web框架日志`), echo.KVxOptX[Other, any](Other{`mock`}))

func Register(name string, title string, otherNames ...string) {
	Categories.Add(name, title, echo.KVxOptX[Other, any](otherNames))
}

func Append(name string, otherNames ...string) {
	item := Categories.GetItem(name)
	if item == nil {
		item = &echo.KVx[Other, any]{
			K: name,
		}
		Categories.AddItem(item)
	}
	item.X = append(item.X, otherNames...)
}
