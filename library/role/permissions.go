package role

import (
	"github.com/coscms/webcore/dbschema"
)

type UserRoleWithPermissions struct {
	*dbschema.NgingUserRole
	Permissions []*dbschema.NgingUserRolePermission `db:"-,relation=role_id:id"`
}

func (u *UserRoleWithPermissions) GetPermissions() []PermissionConfiger {
	r := make([]PermissionConfiger, len(u.Permissions))
	for k, v := range u.Permissions {
		r[k] = v
	}
	return r
}
