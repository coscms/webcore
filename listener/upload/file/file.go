package file

import (
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/db/lib/factory"
)

func init() {
	// - nging_file
	dbschema.DBI.On(factory.EventCreated, func(m factory.Model, _ ...string) (err error) {
		fm := m.(*dbschema.NgingFile)
		if fm.OwnerType == `user` && fm.OwnerId > 0 {
			userM := model.NewUser(fm.Context())
			err = userM.IncrFileSizeAndNum(uint(fm.OwnerId), fm.Size, 1)
		}
		return
	}, `nging_file`)
}
