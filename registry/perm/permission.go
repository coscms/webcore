package perm

import (
	"github.com/coscms/webcore/library/role"
)

func New() *role.RolePermission {
	return role.NewRolePermission()
}

type RolePermission = role.RolePermission
