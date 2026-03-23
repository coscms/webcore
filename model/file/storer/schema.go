package storer

import (
	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

const (
	StorerInfoKey = `NgingStorer`
)

func NewInfo() *Info {
	return &Info{}
}

type Info struct {
	Name  string `json:"name" xml:"name"`
	ID    string `json:"id" xml:"id"`
	cloud *dbschema.NgingCloudStorage
}

func (s *Info) FromStore(v echo.H) *Info {
	s.Name = v.String("name")
	s.ID = v.String("id")
	if s.ID == `0` {
		s.ID = ``
	}
	return s
}

func (s *Info) Cloud(ctx echo.Context, forces ...bool) (*dbschema.NgingCloudStorage, error) {
	var force bool
	if len(forces) > 0 {
		force = forces[0]
	}
	if !force && s.cloud != nil {
		return s.cloud, nil
	}
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	cloudM := dbschema.NewNgingCloudStorage(ctx)
	s.cloud = cloudM
	if len(s.ID) > 0 {
		id := param.AsUint(s.ID)
		err := cloudM.Get(nil, `id`, id)
		return s.cloud, err
	}
	return s.cloud, nil
}
