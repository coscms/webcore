package yaegiplugin

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webcore/library/plugin/yaegiplugin/lib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/webx-top/echo"
)

func Load(files ...string) (err error) {
	if len(files) == 0 {
		pluginGlob := filepath.Join(echo.Wd(), `plugins`) + echo.FilePathSeparator + `*.go`
		files, err = filepath.Glob(pluginGlob)
		if err != nil {
			return
		}
	}
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(lib.Symbols)
	for _, file := range files {
		_, err = i.EvalPath(file)
		if err != nil {
			return
		}
		var v reflect.Value
		v, err = i.Eval("main.Module")
		if err != nil {
			if isUndefined(err) {
				moduleName := filepath.Base(filepath.Dir(file))
				v, err = i.Eval(moduleName + ".Module")
			}
			if err != nil {
				return
			}
		}
		module.Register(v.Interface().(*module.Module))
	}
	return
}

func isUndefined(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.SplitN(err.Error(), `: `, 2)[1] // ./_testdata/main.go:1:28: undefined: main
	return strings.HasPrefix(errMsg, `undefined:`) || strings.HasPrefix(errMsg, `undefined selector:`)
}
