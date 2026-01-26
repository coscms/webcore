package pagebuilder

import (
	"html/template"

	"github.com/coscms/tables"
	"github.com/webx-top/echo"
)

// NewTable 创建并返回一个新的 Table 实例
func NewTable() *Table {
	return &Table{
		Table: tables.New(),
	}
}

type Table struct {
	*tables.Table
}

func (f *Table) RenderContextWithData(ctx echo.Context, data interface{}) template.HTML {
	return f.Table.Render()
}
