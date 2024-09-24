package module

import (
	"path"
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/cmder"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webcore/library/cron"
	"github.com/coscms/webcore/library/nlog"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webcore/registry/settings"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/middleware"
)

type IModule interface {
	Apply()
	Version() float64
}

var _ IModule = &Module{}

type Module struct {
	Startup       string                      // 默认启动项(多个用半角逗号“,”隔开)
	Navigate      func(*Navigate)             // 注册导航菜单
	Extend        map[string]extend.Initer    // 注册扩展配置项
	Cmder         map[string]cmder.Cmder      // 注册命令
	TemplatePath  map[string]string           // 注册模板路径
	AssetsPath    []string                    // 注册素材路径
	SQLCollection func(*config.SQLCollection) // 注册SQL语句
	Dashboard     func(*Dashboard)            // 注册控制面板首页区块
	Route         func(*Router)               // 注册网址路由
	LogParser     map[string]nlog.LogParser   // 注册日志解析器
	Settings      []*settings.SettingForm     // 注册配置选项
	CronJobs      []*cron.Jobx                // 注册定时任务
	DBSchemaVer   float64                     // 设置数据库结构版本号
}

func (m *Module) setNavigate(nc *Navigate) {
	if m.Navigate == nil {
		return
	}
	m.Navigate(nc)
}

func (m *Module) setConfig(*config.Config) {
	if m.Extend == nil {
		return
	}
	for k, v := range m.Extend {
		extend.Register(k, v)
	}
}

func (m *Module) setCmder(*config.CLIConfig) {
	if m.Cmder == nil {
		return
	}
	for k, v := range m.Cmder {
		cmder.Register(k, v)
	}
}

func (m *Module) setTemplate(backendPa *ntemplate.PathAliases, frontendPa *ntemplate.PathAliases) {
	if m.TemplatePath == nil {
		return
	}
	for k, v := range m.TemplatePath {
		SetTemplate(backendPa, frontendPa, k, v)
	}
}

func (m *Module) setAssets(backendSo *middleware.StaticOptions, frontendSo *middleware.StaticOptions) {
	for _, v := range m.AssetsPath {
		SetAssets(backendSo, frontendSo, v)
	}
}

func (m *Module) setSQL(sc *config.SQLCollection) {
	if m.SQLCollection == nil {
		return
	}
	m.SQLCollection(sc)
}

func (m *Module) setDashboard(dd *Dashboard) {
	if m.Dashboard == nil {
		return
	}
	m.Dashboard(dd)
}

func (m *Module) setRoute(r *Router) {
	if m.Route == nil {
		return
	}
	m.Route(r)
}

func (m *Module) setLogParser(parsers map[string]nlog.LogParser) {
	if m.LogParser == nil {
		return
	}
	for k, p := range m.LogParser {
		parsers[k] = p
	}
}

func (m *Module) setSettings() {
	settings.Register(m.Settings...)
}

func (m *Module) setCronJob() {
	for _, jobx := range m.CronJobs {
		jobx.Register()
	}
}

func (m *Module) setDefaultStartup() {
	if len(m.Startup) > 0 {
		if len(config.DefaultStartup) > 0 && !strings.HasPrefix(m.Startup, `,`) {
			config.DefaultStartup += `,` + m.Startup
		} else {
			config.DefaultStartup += m.Startup
		}
	}
}

func (m *Module) Version() float64 {
	return m.DBSchemaVer
}

func (m *Module) Apply() {
	m.setNavigate(NewNavigate())
	m.setConfig(config.FromFile())
	m.setCmder(config.FromCLI())
	m.applyTemplateAndAssets()
	//m.setTemplate(bindata.PathAliases)
	//m.setAssets(bindata.StaticOptions)
	m.setSQL(config.GetSQLCollection())
	m.setDashboard(NewDashboard())
	m.setRoute(NewRouter())
	m.setLogParser(nlog.LogParsers)
	m.setSettings()
	m.setDefaultStartup()
	m.setCronJob()
}

func SetTemplate(backendPa *ntemplate.PathAliases, frontendPa *ntemplate.PathAliases, key string, templatePath string) {
	if len(templatePath) == 0 {
		return
	}
	if templatePath[0] != '.' && templatePath[0] != '/' && !strings.HasPrefix(templatePath, `vendor/`) {
		templatePath = NgingPluginDir + `/` + templatePath
	}
	dir := path.Base(templatePath)
	switch dir {
	case `frontend`:
		if frontendPa != nil {
			frontendPa.Add(key, templatePath)
		}
	case `backend`:
		backendPa.Add(key, templatePath)
	case `template`:
		if frontendPa != nil && com.IsDir(templatePath+`/frontend`) {
			frontendPa.AddAllSubdir(templatePath + `/frontend`) // 支持多主题
			//frontendPa.Add(key, templatePath+`/frontend`)
		}
		if com.IsDir(templatePath + `/backend`) {
			backendPa.Add(key, templatePath+`/backend`)
		}
	}
}

func SetAssets(backendSo *middleware.StaticOptions, frontendSo *middleware.StaticOptions, assetsPath string) {
	if len(assetsPath) == 0 {
		return
	}
	if assetsPath[0] != '.' && assetsPath[0] != '/' && !strings.HasPrefix(assetsPath, `vendor/`) {
		assetsPath = NgingPluginDir + `/` + assetsPath
	}
	dir := path.Base(assetsPath)
	switch dir {
	case `frontend`:
		if frontendSo != nil {
			frontendSo.AddFallback(assetsPath)
		}
	case `backend`:
		backendSo.AddFallback(assetsPath)
	case `assets`:
		if frontendSo != nil && com.IsDir(assetsPath+`/frontend`) {
			frontendSo.AddFallback(assetsPath + `/frontend`)
		}
		if com.IsDir(assetsPath + `/backend`) {
			backendSo.AddFallback(assetsPath + `/backend`)
		}
	}
}
