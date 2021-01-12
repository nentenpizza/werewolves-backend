package game

import (
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Character       Character `json:"character,omitempty"`
	ID              string    `json:"id"`
	Update          chan bool `json:"-"`
	Room            *Room     `json:"room"`
	*websocket.Conn `json:"-"`
	sync.Mutex      `json:"-"`
}

func (p *Player) Kill() {
	p.Character.SetHP(0)
	p.Room.Dead[p.ID] = true
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
	Role      string    `json:"role"`
	Character Character `json:"character"`
	Player    *Player   `json:"me"`
	Room      *Room     `json:"room"`
}

func NewPlayerState(player *Player) Update {
	t := reflect.TypeOf(player.Character)
	return Update{
		Role:      t.Elem().Name(),
		Player:    player,
		Room:      player.Room,
		Character: player.Character,
	}
}
