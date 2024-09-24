package model

import (
	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/webx-top/db"
)

func (u *User) SetSession(users ...*dbschema.NgingUser) {
	userCopy := u.ClearPasswordData(users...)
	u.Context().Session().Set(`user`, &userCopy)
}

func (u *User) UnsetSession() {
	u.Context().Session().Delete(`user`)
}

func (u *User) VerifySession(users ...*dbschema.NgingUser) error {
	var user *dbschema.NgingUser
	if len(users) > 0 {
		user = users[0]
	} else {
		user, _ = u.Context().Session().Get(`user`).(*dbschema.NgingUser)
	}
	if user == nil {
		return nerrors.ErrUserNotLoggedIn
	}
	err := u.Get(nil, db.Cond{`id`: user.Id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		u.UnsetSession()
		return nerrors.ErrUserNotFound
	}
	if !sessionguard.Validate(u.Context(), user.LastIp, `user`, uint64(user.Id)) {
		log.Warn(u.Context().T(`用户“%s”的会话环境发生改变，需要重新登录`, user.Username))
		u.UnsetSession()
		return nerrors.ErrUserNotLoggedIn
	}
	if u.NgingUser.Updated != user.Updated {
		u.SetSession()
		u.Context().Internal().Set(`user`, user)
	}
	return nil
}
