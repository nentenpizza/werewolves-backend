package game

import (
	"github.com/gorilla/websocket"
	"reflect"
	"sync"
)

type Player struct {
	Character       Character `json:"role"`
	ID              string    `json:"id"`
	Update          chan bool `json:"-"`
	Room            *Room     `json:"room"`
	*websocket.Conn `json:"-"`
	sync.Mutex      `json:"-"`
}

func NewPlayer(ID string, conn ...*websocket.Conn) *Player {
	var c *websocket.Conn
	if len(conn) > 0 {
		c = conn[0]
	} else {
		c = nil
	}
	return &Player{ID: ID, Update: make(chan bool), Conn: c}
}

type Players map[string]*Player

type Update struct {
	Role        string   `json:"role"`
	ID          string   `json:"id"`
	DeadPlayers []string `json:"dead_players"`
}

func NewPlayerState(player *Player) Update {
	var deadPlayers = make([]string, 0, len(player.Room.Players))
	for _, p := range player.Room.Dead {
		deadPlayers = append(deadPlayers, p.ID)
	}
	t := reflect.TypeOf(player.Character)
	return Update{Role: t.Elem().Name(), ID: player.ID, DeadPlayers: deadPlayers}
}
