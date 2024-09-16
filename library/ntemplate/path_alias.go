package ntemplate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// NewPathAliases 新建模板路径别名分组设置实例
func NewPathAliases() *PathAliases {
	return &PathAliases{
		aliases: map[string][]string{},
	}
}

// PathAliases 模板路径别名分组设置
type PathAliases struct {
	// aliases map中的key一般为模块定义中module.Module.TemplatePath这个map中的key
	// aliases map中的value一般为模块定义中module.Module.TemplatePath这个map中的value(用于指定模板文件夹路径)
	// aliases map中的value是个切片，意味着在不同的模块定义TemplatePath中，允许定义与其它模块中相同的key
	aliases map[string][]string

	// tmplDirs 所有模板文件夹路径
	tmplDirs []string
}

// TmplDirs 所有模板文件夹路径
func (p *PathAliases) TmplDirs() []string {
	return p.tmplDirs
}

// Aliases 获取全部别名
func (p *PathAliases) Aliases() []string {
	aliases := make([]string, len(p.aliases))
	var i int
	for alias := range p.aliases {
		aliases[i] = alias
		i++
	}
	return aliases
}

// Range 遍历全部别名及其模板文件夹路径
func (p *PathAliases) Range(fn func(string, string) error) (err error) {
	for alias, templateDirs := range p.aliases {
		for _, templateDir := range templateDirs {
			err = fn(alias, templateDir)
			if err != nil {
				return
			}
		}
	}
	return
}

// AddAllSubdir 添加某个路径下以子文件夹为别名的路径(用于添加各个子文件夹为不同主题模板的情况)
func (p *PathAliases) AddAllSubdir(absPath string) error {
	fp, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	dirs, err := fp.ReadDir(-1)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), `.`) {
			continue
		}
		p.Add(dir.Name(), absPath)
	}
	return nil
}

// Add 添加别名和路径
func (p *PathAliases) Add(alias, absPath string) *PathAliases {
	var err error
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		panic(err)
	}
	if !com.InSlice(absPath, p.tmplDirs) {
		p.tmplDirs = append(p.tmplDirs, absPath)
	}
	if p.aliases == nil {
		p.aliases = map[string][]string{}
	}
	if _, ok := p.aliases[alias]; !ok {
		p.aliases[alias] = []string{}
	}
	if !strings.HasSuffix(absPath, echo.FilePathSeparator) {
		absPath += echo.FilePathSeparator
	}
	p.aliases[alias] = append(p.aliases[alias], absPath)
	return p
}

// ParsePrefix 根据模板路径子文件夹为别名解析模板路径为实际存在的真实路径
func (p *PathAliases) ParsePrefix(withAliasPrefixPath string) string {
	rpath, _ := p.ParsePrefixOk(withAliasPrefixPath)
	return rpath
}

// ParsePrefixOk 根据模板路径子文件夹为别名解析模板路径为实际存在的真实路径
func (p *PathAliases) ParsePrefixOk(withAliasPrefixPath string) (string, bool) {
	if len(withAliasPrefixPath) < 3 {
		return withAliasPrefixPath, false
	}
	if withAliasPrefixPath[0] == '/' || withAliasPrefixPath[0] == '.' {
		fi, err := os.Stat(withAliasPrefixPath)
		if err == nil && !fi.IsDir() {
			return withAliasPrefixPath, false
		}
		withAliasPrefixPath = withAliasPrefixPath[1:]
	}
	parts := strings.SplitN(withAliasPrefixPath, `/`, 2)
	if len(parts) != 2 {
		return withAliasPrefixPath, false
	}
	alias := parts[0]
	if opaths, ok := p.aliases[alias]; ok {
		if len(opaths) == 1 {
			return filepath.Join(opaths[0], withAliasPrefixPath), true
		}
		for _, opath := range opaths {
			_tmpl := filepath.Join(opath, withAliasPrefixPath)
			fi, err := os.Stat(_tmpl)
			if err == nil && !fi.IsDir() {
				return _tmpl, true
			}
		}
	}
	return withAliasPrefixPath, false
}

// RestorePrefix 还原真实路径为模板初始路径
func (p *PathAliases) RestorePrefix(fullpath string) string {
	rpath, _ := p.RestorePrefixOk(fullpath)
	return rpath
}

// RestorePrefixOk 还原真实路径为模板初始路径
func (p *PathAliases) RestorePrefixOk(fullpath string) (string, bool) {
	for _, absPaths := range p.aliases {
		for _, absPath := range absPaths {
			if strings.HasPrefix(fullpath, absPath) {
				return filepath.ToSlash(fullpath[len(absPath):]), true
			}
		}
	}
	return fullpath, false
}

// Parse 根据中括号别名前缀解析模板路径为实际存在的真实路径
func (p *PathAliases) Parse(withAliasTagPath string) string {
	rpath, _ := p.ParseOk(withAliasTagPath)
	return rpath
}

// ParseOk 根据中括号别名前缀解析模板路径为实际存在的真实路径
func (p *PathAliases) ParseOk(withAliasTagPath string) (string, bool) {
	if len(withAliasTagPath) < 3 || withAliasTagPath[0] != '[' {
		return withAliasTagPath, false
	}
	withAliasTagPath = withAliasTagPath[1:]
	parts := strings.SplitN(withAliasTagPath, `]`, 2)
	if len(parts) != 2 {
		return withAliasTagPath, false
	}
	alias := parts[0]
	rpath := parts[1]
	if opaths, ok := p.aliases[alias]; ok {
		if len(opaths) == 1 {
			return filepath.Join(opaths[0], rpath), true
		}
		for _, opath := range opaths {
			_tmpl := filepath.Join(opath, rpath)
			fi, err := os.Stat(_tmpl)
			if err == nil && !fi.IsDir() {
				return _tmpl, true
			}
		}
	}
	return rpath, false
}

// Restore 还原真实路径为带中括号别名前缀的模板初始路径
func (p *PathAliases) Restore(fullpath string) string {
	rpath, _ := p.RestoreOk(fullpath)
	return rpath
}

// RestoreOk 还原真实路径为带中括号别名前缀的模板初始路径
func (p *PathAliases) RestoreOk(fullpath string) (string, bool) {
	for alias, absPaths := range p.aliases {
		for _, absPath := range absPaths {
			if strings.HasPrefix(fullpath, absPath) {
				return `[` + alias + `]` + filepath.ToSlash(fullpath[len(absPath):]), true
			}
		}
	}
	return fullpath, false
}
