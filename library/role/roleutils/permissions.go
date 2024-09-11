package roleutils

import (
	"strings"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/role"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type UserRoleWithPermissions struct {
	*dbschema.NgingUserRole
	Permissions []*dbschema.NgingUserRolePermission `db:"-,relation=role_id:id"`
}

func UserRoles(ctx echo.Context) []*role.UserRoleWithPermissions {
	roleList, ok := ctx.Internal().Get(`userRoles`).([]*role.UserRoleWithPermissions)
	if ok {
		return roleList
	}
	roleList = GetRoleList(ctx)
	if len(roleList) > 0 {
		ctx.Internal().Set(`userRoles`, roleList)
	}
	return roleList
}

func GetRoleList(c echo.Context) (roleList []*role.UserRoleWithPermissions) {
	user := backend.User(c)
	if user == nil {
		return nil
	}
	roleM := model.NewUserRole(c)
	if len(user.RoleIds) > 0 {
		roleM.ListByOffset(&roleList, nil, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`id`: db.In(strings.Split(user.RoleIds, `,`))},
		))
	}
	return
}
