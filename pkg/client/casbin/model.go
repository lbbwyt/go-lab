package casbin

import "github.com/casbin/casbin/v2/model"

var Model model.Model

//策略示例：
//p, admin, domain1, data1, read //domain1域下，角色admin拥有data1的读权限
//g, lbb, admin, domain1  //用户lbb 拥有domain1域下的 admin角色

//用户 : userID (用户对应存储的为用户唯一ID)
//域： domainID（域对应存储的为团队唯一ID）
//资源：存储的为资源唯一ID
//动作：存储的为动作标识 “Module：OpLevel：OpClass：OpName”

func init() {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")
	m.AddDef("p", "p", "sub, dom, obj, act")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub, r.dom) && keyMatch(r.dom, p.dom) && keyMatch(r.obj, p.obj) && keyMatch(r.act, p.act)")
	//m.AddDef("m", "m", "g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act")
	Model = m
}
