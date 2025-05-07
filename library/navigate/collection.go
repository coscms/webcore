package navigate

type NavigateType string

const (
	Left   NavigateType = `left`
	Top    NavigateType = `top`
	Right  NavigateType = `right`
	Bottom NavigateType = `bottom`
)

func NewProjectNavigates(kind, baseProject string) *ProjectNavigates {
	return &ProjectNavigates{
		Navigates:        &Navigates{},
		baseProject:      baseProject,
		kind:             kind,
		projectNavigates: map[string]*Navigates{},
		Projects:         NewProjects(),
	}
}

type ProjectNavigates struct {
	*Navigates
	baseProject      string
	kind             string
	projectNavigates map[string]*Navigates
	Projects         *Projects
}

func (p *ProjectNavigates) Init() {
	p.Projects.InitURLsIdent()
	p.GroupedChildren(p.kind)
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

// func (p *ProjectNavigates) Add(typ NavigateType, nav *List) {
// 	p.navigates.Add(typ, nav)
// }

// func (p *ProjectNavigates) AddItems(typ NavigateType, index int, items ...*Item) {
// 	p.navigates.AddItems(typ, index, items...)
// }

// func (p *ProjectNavigates) AddTopItems(index int, items ...*Item) {
// 	p.navigates.AddTopItems(index, items...)
// }

// func (p *ProjectNavigates) AddLeftItems(index int, items ...*Item) {
// 	p.navigates.AddLeftItems(index, items...)
// }

// func (p *ProjectNavigates) AddRightItems(index int, items ...*Item) {
// 	p.navigates.AddRightItems(index, items...)
// }

// func (p *ProjectNavigates) AddBottomItems(index int, items ...*Item) {
// 	p.navigates.AddBottomItems(typ)
// }

// func (p *ProjectNavigates) Get(typ NavigateType) (nav *List) {
// 	return p.navigates.Get(typ)
// }

// func (p *ProjectNavigates) GetTop() *List {
// 	return p.navigates.GetTop()
// }

// func (p *ProjectNavigates) GetLeft() *List {
// 	return p.navigates.GetLeft()
// }

// func (p *ProjectNavigates) GetRight() *List {
// 	return p.navigates.GetRight()
// }

// func (p *ProjectNavigates) GetBottom() *List {
// 	return p.navigates.GetBottom()
// }

// func (p *ProjectNavigates) Remove(typ NavigateType) bool {
// 	return p.navigates.Remove(typ)
// }

type Navigates map[NavigateType]*List

func (n *Navigates) Add(typ NavigateType, nav *List) {
	(*n)[typ] = nav
}

func (n *Navigates) GroupedChildren(prefix string) {
	for nType, nRows := range *n {
		nRows.groupedChildren(prefix + `.` + string(nType))
	}
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
