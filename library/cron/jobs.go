package cron

import (
	"sort"

	"github.com/webx-top/echo"
)

var (
	// systemJobs 系统Job
	systemJobs = map[string]Jobx{}
)

func ListSystemJobs() echo.KVList {
	kvList := echo.KVList{}
	names := make([]string, 0, len(systemJobs))
	for name := range systemJobs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		sj := systemJobs[name]
		kvList = append(kvList, echo.NewKV(name, sj.Description).SetHKV(`example`, sj.Example))
	}
	return kvList
}

func Register(name string, fn RunnerGetter, example string, description string) {
	AddSystemJob(name, fn, example, description)
}

// AddSystemJob 添加系统Job
func AddSystemJob(name string, fn RunnerGetter, example string, description string) {
	systemJobs[name] = Jobx{
		Name:         name,
		Example:      example,
		Description:  description,
		RunnerGetter: fn,
	}
}

type Jobx struct {
	Name         string
	Example      string //">funcName:param"
	Description  string
	RunnerGetter RunnerGetter
}

func (j *Jobx) Register() {
	systemJobs[j.Name] = *j
}
