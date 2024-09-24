package file_test

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/coscms/webcore/library/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	myTesting "github.com/webx-top/echo/testing"
	"github.com/webx-top/echo/testing/test"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/testutils"
	_ "github.com/coscms/webcore/listener/upload/file"
	modelFile "github.com/coscms/webcore/model/file"
)

// TestUpdateEmbedded 测试图片修改
func TestUpdateEmbedded(t *testing.T) {
	ownerID := uint64(1)

	testutils.InitConfig()

	e := echo.New()
	e.Use(httpserver.Transaction())
	req, resp := myTesting.NewRequestAndResponse(`GET`, `/`)
	ctx := e.NewContext(req, resp)
	ctx.SetTransaction(factory.NewParam())
	tables := []string{}
	for table := range dbschema.DBI.Events {
		tables = append(tables, table)
	}
	echo.Dump(tables)
	test.Contains(t, tables, `nging_user`)
	userM := dbschema.NewNgingUser(ctx)
	assert.True(t, dbschema.DBI.Events.Exists(factory.EventUpdated, userM))
	fileM := modelFile.NewFile(ctx)
	userM.Get(nil, `id`, ownerID)
	if len(userM.Avatar) > 0 {
		userM.Avatar = ``
	} else {
		userM.Avatar = `/public/upload/user/` + fmt.Sprint(ownerID) + `/avatar.jpg`
		err := fileM.GetByViewURL(userM.Avatar)
		if err != nil {
			if err != db.ErrNoMoreRows {
				panic(err)
			}
			fileM.SaveName = path.Base(userM.Avatar)
			fileM.ViewUrl = userM.Avatar
			fileM.SavePath = strings.TrimPrefix(userM.Avatar, `/`)
			fileM.Subdir = `avatar`
			fileM.Type = `image`
			fileM.Size = 100
			fileM.Ext = path.Ext(fileM.SaveName)
			fileM.OwnerType = `user`
			fileM.OwnerId = ownerID
			fileM.Mime = `image/jpeg`
			fileM.StorerName = `local`
			_, err = fileM.Insert()
			assert.NoError(t, err)
		}
	}
	if err := userM.Update(nil, `id`, ownerID); err != nil {
		panic(err)
	}
	err := fileM.Get(nil, db.And(
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: `user`},
		db.Cond{`subdir`: `avatar`},
	))
	if err != nil {
		panic(err)
	}
	if len(userM.Avatar) > 0 {
		err = fileM.GetByViewURL(userM.Avatar)
		if err != nil {
			panic(err)
		}
	}
	em := modelFile.NewEmbedded(ctx, fileM)
	num, err := em.Count(nil, db.And(
		db.Cond{`table_id`: ownerID},
		db.Cond{`table_name`: `nging_user`},
		db.Cond{`field_name`: `avatar`},
	))
	if err != nil && err != db.ErrNoMoreRows {
		panic(err)
	}
	if len(userM.Avatar) == 0 { // 头像删除
		fmt.Println(`头像删除`)
		test.Eq(t, uint(0), fileM.UsedTimes)
		test.Eq(t, int64(0), num)
	} else { // 添加头像
		fmt.Println(`添加头像`)
		test.Eq(t, uint(1), fileM.UsedTimes)
		test.Eq(t, int64(1), num)
	}
}
