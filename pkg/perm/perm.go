package perm

import "errors"

const (
	RoleDel    = 1 << 0
	RoleUpdate = 1 << 1
	RoleCreate = 1 << 2
	//...
	RoleAuth = 1 << 63
)

var roleId2NameMap = map[uint64]string{
	RoleDel:    "删除",
	RoleUpdate: "更新",
	RoleCreate: "创建",
	RoleAuth:   "授权角色",
}

//模拟数据库存储， 用于存储用户分配的权限
var userRoleMap = map[string]map[string]uint64{
	"452700387": {
		"A": RoleCreate | RoleUpdate, //定义用户452700387对资源A有创建和修改的权限
	},

	"452700388": {
		"*": RoleAuth, //定义用户452700388对任何资源有授权的权限
	},
}

type Perm struct {
}

//获取用户的权限
func (p *Perm) GetUserPerm(userId string) map[string]uint64 {
	if item, ok := userRoleMap[userId]; ok {
		return item
	}
	return map[string]uint64{}
}

// 授权人授权对某资源的操作权限给用户
func (p *Perm) Auth(authorizer string, resource string, toUser string, role uint64) error {
	//获取授权人的权限
	hasPermission := p.HasPermission(authorizer, resource, RoleAuth)
	if !hasPermission {
		return errors.New("has no auth permission")
	}

	if item, ok := userRoleMap[toUser]; ok {
		if rv, rok := item[resource]; rok {
			item[resource] = role
		} else {
			if !(rv&role == role) {
				// 未授权 该权限时，才授权， 否则就不用授权了
				item[resource] = rv + role
			}
		}
		userRoleMap[toUser] = item

	} else {
		userRoleMap[toUser] = map[string]uint64{
			resource: role,
		}
	}
	return nil
}

//用户鉴权
func (p *Perm) HasPermission(userId, resource string, role uint64) bool {
	perms := p.GetUserPerm(userId)
	hasPermission := false

	for k, v := range perms {
		if k != "*" && k != resource {
			continue
		}
		//判断是否对该资源的授权权限
		if v&role == role {
			hasPermission = true
			break
		}
	}
	return hasPermission
}

//回收权限
func (p *Perm) RevokeAuth(authorizer string, resource string, fromUser string, role uint64) error {
	//获取授权人的权限
	hasPermission := p.HasPermission(authorizer, resource, RoleAuth)
	if !hasPermission {
		return errors.New("has no auth permission")
	}

	//已存在授权时回收
	if item, ok := userRoleMap[fromUser]; ok {
		if rv, rok := item[resource]; rok {
			if rv&role == role {
				// 已授权时才取消
				item[resource] = rv - role
			}
		}
		userRoleMap[fromUser] = item
	}
	return nil
}
