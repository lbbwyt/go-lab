package perm

import (
	"fmt"
	"testing"
)

func TestPerm_Auth(t *testing.T) {

	p := new(Perm)
	p.Auth("452700388", "B", "452700399", RoleCreate)
	fmt.Println(p.HasPermission("452700399", "B", RoleCreate))

	fmt.Println(userRoleMap)

	p.RevokeAuth("452700388", "B", "452700399", RoleCreate)

	fmt.Println(p.HasPermission("452700399", "B", RoleCreate))
	fmt.Println(userRoleMap)

	p.Auth("452700388", "B", "452700399", RoleCreate)
	fmt.Println(p.HasPermission("452700399", "B", RoleCreate))

	fmt.Println(userRoleMap)
}
