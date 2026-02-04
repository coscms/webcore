package pagebuilder

import (
	"html/template"

	"github.com/coscms/tables"
	"github.com/webx-top/echo"
)

var _ echo.RenderContextWithData = (*Table)(nil)

// NewTable 创建并返回一个新的 Table 实例
func NewTable() *Table {
	return &Table{
		Table: tables.New(),
	}
}

type Table struct {
	*tables.Table
}

func (f *Table) RenderWithData(ctx echo.Context, data interface{}) template.HTML {
	return f.Table.Render()
}
