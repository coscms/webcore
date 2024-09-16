package dashboard

import (
	"html/template"
	"strings"

	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

func NewPage(key string, atmpls ...map[string][]string) *Page {
	var tmpls map[string][]string
	if len(atmpls) > 0 {
		tmpls = atmpls[0]
	}
	if tmpls == nil {
		tmpls = map[string][]string{}
	}
	return &Page{
		Key:   key,
		Tmpls: tmpls,
		hooks: map[string][]func(echo.Context) error{},
	}
}

type Page struct {
	Key   string
	Tmpls map[string][]string
	hooks map[string][]func(echo.Context) error
}

func (s *Page) AddTmpl(position string, tmpl ...string) *Page {
	if _, ok := s.Tmpls[position]; !ok {
		s.Tmpls[position] = []string{}
	}
	s.Tmpls[position] = append(s.Tmpls[position], tmpl...)
	return s
}

func (s *Page) Tmpl(position string) []string {
	return s.Tmpls[position]
}

func (s *Page) On(method string, hook func(echo.Context) error) *Page {
	method = strings.ToUpper(method)
	if _, ok := s.hooks[method]; !ok {
		s.hooks[method] = []func(echo.Context) error{}
	}
	s.hooks[method] = append(s.hooks[method], hook)
	return s
}

func (s *Page) Fire(ctx echo.Context) error {
	method := strings.ToUpper(ctx.Method())
	hooks, ok := s.hooks[method]
	if !ok {
		return nil
	}
	errs := common.NewErrors()
	for _, hook := range hooks {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}

func Render(ctx echo.Context, tmpl string, data interface{}) template.HTML {
	b, err := ctx.Fetch(tmpl, data)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(engine.Bytes2str(b))
}
