package model

type MonsterParam string

type Monster struct {
	Name string
}

func NewMonster() Monster {
	return Monster{Name: "kitty"}
}
func NewMonsterWithParam(name MonsterParam) Monster {
	return Monster{
		Name: string(name),
	}
}
