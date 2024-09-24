package navigate

type NavigateType string

const (
	Left   NavigateType = `left`
	Top    NavigateType = `top`
	Right  NavigateType = `right`
	Bottom NavigateType = `bottom`
)

func NewProjectNavigates(baseProject string) *ProjectNavigates {
	return &ProjectNavigates{
		Navigates:        &Navigates{},
		baseProject:      baseProject,
		projectNavigates: map[string]*Navigates{},
		Projects:         NewProjects(),
	}
}

type ProjectNavigates struct {
	*Navigates
	baseProject      string
	projectNavigates map[string]*Navigates
	Projects         *Projects
}

func (p *ProjectNavigates) AddProject(index int, list ...*ProjectItem) {
	p.Projects.Add(index, list...)
	for _, item := range list {
		p.Project(item.Ident).Add(Left, item.NavList)
	}
}

func (p *ProjectNavigates) AddNavList(name string, ident string, url string, navList *List) {
	proj := p.Projects.Get(ident)
	if proj == nil {
		p.AddProject(-1, NewProject(name, ident, url, navList))
		return
	}
	proj.NavList.Add(-1, *navList...)
}

func (p *ProjectNavigates) Project(project string) *Navigates {
	if p.baseProject == project {
		return p.Navigates
	}
	nav, ok := p.projectNavigates[project]
	if !ok {
		nav = &Navigates{}
		p.projectNavigates[project] = nav
	}
	return nav
}

func (p *ProjectNavigates) RemoveProject(project string) {
	delete(p.projectNavigates, project)
}

type Navigates map[NavigateType]*List

func (n *Navigates) Add(typ NavigateType, nav *List) {
	(*n)[typ] = nav
}

func (n *Navigates) AddItems(typ NavigateType, index int, items ...*Item) {
	nav := n.Get(typ)
	if nav == nil {
		nav = &List{}
		(*n)[typ] = nav
	}
	nav.Add(index, items...)
}

func (n *Navigates) AddTopItems(index int, items ...*Item) {
	n.AddItems(Top, index, items...)
}

func (n *Navigates) AddLeftItems(index int, items ...*Item) {
	n.AddItems(Left, index, items...)
}

func (n *Navigates) AddRightItems(index int, items ...*Item) {
	n.AddItems(Right, index, items...)
}

func (n *Navigates) AddBottomItems(index int, items ...*Item) {
	n.AddItems(Bottom, index, items...)
}

func (n *Navigates) Get(typ NavigateType) (nav *List) {
	nav = (*n)[typ]
	return
}

func (n *Navigates) GetTop() *List {
	return n.Get(Top)
}

func (n *Navigates) GetLeft() *List {
	return n.Get(Left)
}

func (n *Navigates) GetRight() *List {
	return n.Get(Right)
}

func (n *Navigates) GetBottom() *List {
	return n.Get(Bottom)
}

func (n *Navigates) Remove(typ NavigateType) bool {
	_, ok := (*n)[typ]
	if ok {
		delete(*n, typ)
	}
	return ok
}
