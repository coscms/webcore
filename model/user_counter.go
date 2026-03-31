package model

import (
	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func (u *User) IncrLoginFails() error {
	return u.NgingUser.UpdateField(nil, `login_fails`, db.Raw(`login_fails+1`), `id`, u.Id)
}

func (u *User) ResetLoginFails() error {
	return u.NgingUser.UpdateField(nil, `login_fails`, 0, `id`, u.Id)
}

func (u *User) IncrFileSizeAndNum(userID uint, fileSize uint64, fileNum uint64) error {
	return u.NgingUser.UpdateFields(nil, echo.H{
		`file_size`: db.Raw(`file_size+` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num+` + param.AsString(fileNum)),
	}, `id`, userID)
}

func (u *User) DecrFileSizeAndNum(userID uint, fileSize uint64, fileNum uint64) error {
	return u.NgingUser.UpdateFields(nil, echo.H{
		`file_size`: db.Raw(`file_size-` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num-` + param.AsString(fileNum)),
	}, `id`, userID)
}

func (u *User) SafeDecrFileSizeAndNum(userID uint, fileSize uint64, fileNum uint64) error {
	return u.NgingUser.UpdateFields(nil, map[string]interface{}{
		`file_size`: db.Raw(`file_size-` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num-` + param.AsString(fileNum)),
	}, db.And(
		db.Cond{`id`: userID},
		db.Cond{`file_size`: db.Gte(fileSize)},
		db.Cond{`file_num`: db.Gte(fileNum)},
	))
}

// RecountFile 重新统计客户上传的文件数量和尺寸
func (u *User) RecountFile(userID ...uint) (totalNum uint64, totalSize uint64, err error) {
	ownerID := u.Id
	if len(userID) > 0 && userID[0] > 0 {
		ownerID = userID[0]
	}
	fileM := dbschema.NewNgingFile(u.Context())
	recv := echo.H{}
	err = fileM.NewParam().SetMW(func(r db.Result) db.Result {
		return r.Select(db.Raw(`SUM(size) AS c`), db.Raw(`COUNT(1) AS n`))
	}).SetRecv(&recv).SetArgs(db.And(
		db.Cond{`owner_type`: `user`},
		db.Cond{`owner_id`: ownerID},
	)).One()
	if err != nil {
		return
	}
	totalNum = recv.Uint64(`n`)
	totalSize = recv.Uint64(`c`)
	err = u.NgingUser.UpdateFields(nil, map[string]interface{}{
		`file_size`: totalSize,
		`file_num`:  totalNum,
	}, db.Cond{`id`: ownerID})
	return
}
