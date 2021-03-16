package model

type PlayerParam string

type Player struct {
	Name string
}

func NewPlayer(name string) Player {
	return Player{Name: name}
}

func NewPlayerWithParam(name PlayerParam) Player {
	return Player{Name: string(name)}
}
