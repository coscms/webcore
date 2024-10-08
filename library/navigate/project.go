package navigate

import "strings"

func NewProjects() *Projects {
	return &Projects{
		List: &ProjectList{},
		Hash: map[string]*ProjectItem{},
	}
}

func NewProject(name string, ident string, url string, navLists ...*List) *ProjectItem {
	var navList *List
	if len(navLists) > 0 {
		navList = navLists[0]
	}
	if navList == nil {
		navList = &List{}
	}
	return &ProjectItem{
		Name:    name,
		Ident:   ident,
		URL:     url,
		NavList: navList,
	}
}

type Projects struct {
	urlsIdent map[string]string //网址路由=>项目标识(Ident)
	List      *ProjectList
	Hash      map[string]*ProjectItem //项目标识(Ident)=>项目信息
}

func (p *Projects) URLsIdent() map[string]string {
	if p.urlsIdent != nil {
		return p.urlsIdent
	}
	return p.InitURLsIdent().urlsIdent
}

func (p *Projects) First(notEmptyOpts ...bool) *ProjectItem {
	var notEmpty bool
	if len(notEmptyOpts) > 0 {
		notEmpty = notEmptyOpts[0]
	}
	if p.List != nil && len(*p.List) > 0 {
		if !notEmpty {
			return (*p.List)[0]
		}
		for _, item := range *p.List {
			if item != nil {
				return item
			}
		}
	}
	return nil
}

func (p *Projects) InitURLsIdent() *Projects {
	p.urlsIdent = map[string]string{}
	for ident, proj := range p.Hash {
		if proj.NavList == nil {
			continue
		}
		for _, urlPath := range proj.NavList.FullPath(``) {
			p.urlsIdent[urlPath] = ident
		}
	}
	return p
}

func (p *Projects) Get(ident string) *ProjectItem {
	if item, ok := p.Hash[ident]; ok {
		return item
	}
	return nil
}
func (p *Projects) Remove(index int) *Projects {
	if len(*p.List) <= index {
		return p
	}
	ident := (*p.List)[index].Ident
	p.List.Remove(index)
	delete(p.Hash, ident)
	return p
}
func (p *Projects) Add(index int, list ...*ProjectItem) *Projects {
	for _, item := range list {
		ident := item.Ident
		if _, ok := p.Hash[ident]; ok {
			panic(`Project already exists: ` + item.Ident)
		}
		p.Hash[ident] = item
	}
	p.List.Add(index, list...)
	return p
}
func (p *Projects) Set(index int, list ...*ProjectItem) *Projects {
	p.List.Set(index, list...)
	for _, item := range list {
		p.Hash[item.Ident] = item
	}
	return p
}

func (p *Projects) GetIdent(urlPath string) string {
	urlPath = strings.TrimPrefix(urlPath, `/`)
	if ident, ok := p.URLsIdent()[urlPath]; ok {
		return ident
	}
	return ``
}

func (p *Projects) RemoveByIdent(ident string) {
	index := p.List.SearchIdent(ident)
	if index < 0 {
		return
	}
	p.Remove(index)
}

type ProjectList []*ProjectItem

type ProjectItem struct {
	Name    string
	Ident   string
	URL     string
	NavList *List
}

func (a *ProjectItem) Is(ident string) bool {
	return a.Ident == ident
}

func (a *ProjectItem) GetName() string {
	return a.Name
}

func (a *ProjectItem) GetIdent() string {
	return a.Ident
}

func (a *ProjectItem) GetURL() string {
	return a.URL
}

func (a *ProjectList) SearchIdent(ident string) int {
	r := -1
	for key, item := range *a {
		if item == nil {
			continue
		}
		if r == -1 {
			r = key
		}
		if len(ident) == 0 {
			return r
		}
		if item.Is(ident) {
			return key
		}
	}
	return r
}

// Remove 删除元素
func (a *ProjectList) Remove(index int) *ProjectList {
	if index < 0 {
		*a = (*a)[0:0]
		return a
	}
	size := len(*a)
	if size > index {
		if size > index+1 {
			*a = append((*a)[0:index], (*a)[index+1:]...)
		} else {
			*a = (*a)[0:index]
		}
	}
	return a
}

// Set 设置元素
func (a *ProjectList) Set(index int, list ...*ProjectItem) *ProjectList {
	if len(list) == 0 {
		return a
	}
	if index < 0 {
		*a = append(*a, list...)
		return a
	}
	size := len(*a)
	if size > index {
		(*a)[index] = list[0]
		if len(list) > 1 {
			a.Set(index+1, list[1:]...)
		}
		return a
	}
	for start, end := size, index-1; start < end; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}

// Add 添加列表项
func (a *ProjectList) Add(index int, list ...*ProjectItem) *ProjectList {
	if len(list) == 0 {
		return a
	}
	if index < 0 {
		*a = append(*a, list...)
		return a
	}
	size := len(*a)
	if size > index {
		list = append(list, (*a)[index])
		(*a)[index] = list[0]
		if len(list) > 1 {
			a.Add(index+1, list[1:]...)
		}
		return a
	}
	for start, end := size, index-1; start < end; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}
