package httpserver

import (
	"net/http"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webcore/library/route"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/validator"
)

type HTTPServers struct {
	Backend  *HTTPServer
	Frontend *HTTPServer
}

func (a *HTTPServers) Clear() {
	a.Backend.Clear()
	a.Frontend.Clear()
}

func New(kind string) *HTTPServer {
	s := &HTTPServer{
		Name:                 kind,
		Dashboard:            dashboard.New(),
		Template:             ntemplate.New(kind, nil, true),
		DefaultStaticRootURL: `/public/`,
		DefaultTemplateDir:   `./template/` + kind,
		DefaultAssetsDir:     `./public/assets/` + kind,
		DefaultAssetsURLPath: `/public/assets/` + kind,
	}
	s.TemplateDir = s.DefaultTemplateDir
	s.AssetsDir = s.DefaultAssetsDir
	s.AssetsURLPath = s.DefaultAssetsURLPath
	s.DefaultAvatarURL = s.AssetsURLPath + `/images/user_128.png`
	s.StaticRootURLPath = s.DefaultStaticRootURL
	return s
}

type HTTPServer struct {
	Name      string
	Navigate  *navigate.ProjectNavigates
	Router    route.IRegister
	Dashboard *dashboard.Dashboard
	Template  *ntemplate.Template
	TmplMgr   driver.Manager

	// for web framework
	StaticOptions         *middleware.StaticOptions
	StaticMW              echo.MiddlewareFunc
	KeepExtensionPrefixes []string
	RouteDefaultExtension string
	DefaultTemplateDir    string // 模板路径默认值
	DefaultAssetsDir      string // 素材路径默认值
	DefaultAssetsURLPath  string // 素材网址路径默认值
	DefaultStaticRootURL  string
	StaticRootURLPath     string
	TemplateDir           string                                   //模板文件夹
	AssetsDir             string                                   //素材文件夹
	AssetsURLPath         string                                   //素材网址路径
	DefaultAvatarURL      string                                   //默认头像网址
	RendererDo            func(driver.Driver)                      //模板引擎配置函数
	TmplCustomParser      func(tmpl string, content []byte) []byte //模板自定义解析函数
	ParseStrings          map[string]string                        //模板内容替换
	ParseStringFuncs      map[string]func() string                 //模板内容替换函数
	Middlewares           []interface{}
	GlobalFuncMap         map[string]interface{}
	FuncSetters           []func(echo.Context) error
	HostCheckerRegexpKey  string
	renderOptions         *render.Config
	language              *language.Language
}

func (h *HTTPServer) Clear() {
	if h.Router != nil {
		h.Router.Clear()
	}
}

func (h *HTTPServer) SetPrefix(prefix string) *HTTPServer {
	h.Router.SetPrefix(prefix)
	h.AssetsURLPath = prefix + h.DefaultAssetsURLPath
	h.DefaultAvatarURL = h.AssetsURLPath + `/images/user_128.png`
	h.StaticRootURLPath = prefix + h.DefaultStaticRootURL
	return h
}

func (h *HTTPServer) Prefix() string {
	return h.Router.Prefix()
}

func (h *HTTPServer) SetRouter(router route.IRegister) *HTTPServer {
	h.Router = router
	return h
}

func (h *HTTPServer) SetNavigate(nav *navigate.ProjectNavigates) *HTTPServer {
	h.Navigate = nav
	return h
}

func (h *HTTPServer) SetTmplCustomParser(parser func(tmpl string, content []byte) []byte) *HTTPServer {
	h.TmplCustomParser = parser
	return h
}

func (h *HTTPServer) SetKeepExtensionPrefixes(keepExtensionPrefixes []string) *HTTPServer {
	h.KeepExtensionPrefixes = keepExtensionPrefixes
	return h
}

func (h *HTTPServer) Renderer() driver.Driver {
	if h.renderOptions == nil {
		return nil
	}
	return h.renderOptions.Renderer()
}

func (h *HTTPServer) I18n() *language.Language {
	return h.language
}

func (h *HTTPServer) SetRenderDataWrapper(dataWrapper echo.DataWrapper) *HTTPServer {
	h.Router.Echo().SetRenderDataWrapper(dataWrapper)
	return h
}

func (h *HTTPServer) GetStaticMW() echo.MiddlewareFunc {
	// 注册静态资源文件(网站素材文件)
	if h.StaticMW == nil && h.StaticOptions != nil {
		if len(h.StaticOptions.Root) == 0 {
			h.StaticOptions.Root = h.AssetsDir
		}
		if len(h.StaticOptions.Path) == 0 {
			h.StaticOptions.Path = h.Prefix() + "/public/assets/" + h.Name
		}
		h.StaticOptions.TrimPrefix = h.Prefix()
		h.StaticMW = middleware.Static(h.StaticOptions)
	}
	return h.StaticMW
}

func (h *HTTPServer) PublicHandler(handler interface{}, meta ...echo.H) echo.Handler {
	return PublicHandler(h.Router, handler, meta...)
}

func (h *HTTPServer) GuestHandler(handler interface{}, meta ...echo.H) echo.Handler {
	return GuestHandler(h.Router, handler, meta...)
}

func (h *HTTPServer) Apply() {
	e := h.Router.Echo()
	//e.SetRenderDataWrapper(echo.DefaultRenderDataWrapper)
	if len(h.Router.Prefix()) > 0 {
		e.Pre(FixedUploadURLPrefix())
	}
	e.Use(middleware.Recover())
	if len(h.HostCheckerRegexpKey) > 0 {
		e.Use(HostChecker(h.HostCheckerRegexpKey))
	}
	e.Use(MaxRequestBodySize)
	if len(h.Middlewares) == 0 {
		if !config.FromFile().Sys.DisableHTTPLog {
			e.Use(middleware.Log())
		}
	} else {
		e.Use(h.Middlewares...)
	}

	// 注册静态资源文件(网站素材文件)
	if staticMW := h.GetStaticMW(); staticMW != nil {
		e.Use(staticMW)
	}

	// 启用session
	e.Use(session.Middleware(config.SessionOptions, config.AutoSecure))
	// 启用多语言支持
	h.language = language.New(&config.FromFile().Language)
	e.Use(h.language.Middleware())

	// 启用Validation
	e.Use(validator.Middleware())

	// 事物支持
	e.Use(Transaction())
	// 注册模板引擎
	h.renderOptions = &render.Config{
		TmplDir: h.TemplateDir,
		Engine:  `standard`,
		ParseStrings: map[string]string{
			`__TMPL__`: h.TemplateDir,
		},
		DefaultHTTPErrorCode: http.StatusOK,
		Reload:               true,
		ErrorPages:           config.FromFile().Sys.ErrorPages,
		ErrorProcessors:      ErrorProcessors,
		FuncMapGlobal:        h.GlobalFuncMap,
		CustomParser:         h.TmplCustomParser,
	}
	for key, val := range h.ParseStrings {
		h.renderOptions.ParseStrings[key] = val
	}
	for key, val := range h.ParseStringFuncs {
		h.renderOptions.ParseStringFuncs[key] = val
	}
	if h.RendererDo != nil {
		h.renderOptions.AddRendererDo(h.RendererDo)
	}
	funcSetters := make([]echo.HandlerFunc, 0, len(h.FuncSetters)+1)
	for _, f := range h.FuncSetters {
		funcSetters = append(funcSetters, f)
	}
	funcSetters = append(funcSetters, ErrorPageFunc)
	h.renderOptions.AddFuncSetter(funcSetters...)

	h.renderOptions.ApplyTo(e, h.TmplMgr)

	if len(h.RouteDefaultExtension) > 0 {
		e.SetDefaultExtension(h.RouteDefaultExtension)
		if len(h.KeepExtensionPrefixes) > 0 {
			e.Pre(TrimPathSuffix(h.KeepExtensionPrefixes...))
		}
	}
}
