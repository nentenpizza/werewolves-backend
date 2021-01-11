package game

import "sync"

type Player struct {
	Role Character `json:"role"`
	ID   string    `json:"id"`
	sync.Mutex
}

func NewPlayer(ID string) *Player {
	return &Player{ID: ID}
}

type Players map[string]*Player
