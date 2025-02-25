package request

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
)

// Login 登录表单值(暂未使用)
type Login struct {
	User string `validate:"required,username"`
	Pass string `validate:"required,min=8,max=64"`
	Code string `validate:"required"`
}

func (r *Login) BeforeValidate(ctx echo.Context) error {
	if len(r.Pass) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入密码`).SetZone(`password`)
	}
	passwd, err := backend.DecryptPassword(ctx, r.User, r.Pass)
	if err != nil {
		err = ctx.NewError(code.InvalidParameter, `密码解密失败: %v`, err).SetZone(`password`)
	} else {
		r.Pass = passwd
	}
	return err
}
