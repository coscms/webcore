package pagebuilder

// <a href="{{.URL}}" class="btn{{if .Float}} pull-{{.Float}}{{end}}{{if .Class}} {{.Class}}{{else}} btn-success{{end}}"{{if .Style}} style="{{.Style}}"{{end}}{{range $k,$v:=.Attrs}} {{$k}}="{{$v}}"{{end}}>
//
//	{{if .Icon}}<i class="fa fa-{{.Icon}}"></i>{{end}}
//	{{.Title|$.T}}
//
// </a>
type Button struct {
	Title string
	URL   string
	Icon  string
	Class string
	Style string
	Attrs map[string]string
	Float string // left / right
}
