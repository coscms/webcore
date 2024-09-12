package ntemplate

import (
	"github.com/webx-top/echo/middleware/render/driver"
)

// NewManager 新建模板文件系统管理驱动实例
func NewManager(mgr driver.Manager, pa PathAliases) driver.Manager {
	return &manager{
		Manager: mgr,
		pa:      pa,
	}
}

// manager 模板文件系统管理
type manager struct {
	driver.Manager
	pa PathAliases
}

// AddCallback 添加"模板文件发生变动"时的回调处理函数
func (m *manager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	originalCb := callback
	callback = func(name, typ, event string) {
		name = m.pa.RestorePrefix(name)
		originalCb(name, typ, event)
	}
	m.Manager.AddCallback(rootDir, callback)
}
