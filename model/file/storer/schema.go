package storer

import (
	"fmt"
	"strconv"

	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
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
	if len(s.ID) > 0 && s.Name == `s3` {
		id, err := StorerIDToNumber(s.ID)
		if err != nil {
			return s.cloud, fmt.Errorf(`[storer.Info]%w`, err)
		}
		err = cloudM.Get(nil, `id`, id)
		return s.cloud, err
	}
	return s.cloud, nil
}

func StorerIDToNumber(storerID string) (uint64, error) {
	if len(storerID) == 0 {
		return 0, nil
	}
	id, err := strconv.ParseUint(storerID, 10, strconv.IntSize)
	if err != nil {
		err = fmt.Errorf(`[storerID] failed to strconv.ParseUint(%q): %w`, storerID, err)
	}
	return id, err
}
