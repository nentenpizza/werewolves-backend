package werewolves

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Role            string    `json:"role"`
	Character       Character `json:"character"`
	Voted           bool      `json:"voted"`
	ID              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Update          chan bool `json:"-"`
	Room            *Room     `json:"room"`
	*websocket.Conn `json:"-"`
	sync.Mutex      `json:"-"`
}

func (p *Player) Vote(pID string) Action {
	p.Lock()
	defer p.Unlock()
	return NewAction(VoteAction, func(r *Room) error {
		if p.Voted == true {
			return errors.New("game: player already voted")
		}
		_, ok := r.Players[pID]
		if !ok {
			return errors.New("game: player is not in room")
		}

		if r.State == DayVoting {
			r.Votes[pID]++ // мне насрать
			p.Voted = true // мне насрать
		} else if r.State == Night { // мне насрать
			if p.Role == "Werewolf" || p.Role == "AplhaWerewolf" { // мне насрать
				r.Votes[pID]++ // мне насрать
				p.Voted = true // мне насрать
			} else {
				return errors.New("game: can not vote in night as long as you not werewolf")
			}
		} else {
			return errors.New("game: can not vote in state discuss")
		}

		return nil
	})
}

func (p *Player) Kill() {
	p.Character.SetHP(0)
	p.Room.Dead[p.ID] = true
}

func NewPlayer(name string, conn ...*websocket.Conn) *Player {
	id := uuid.New().String()
	var c *websocket.Conn
	if len(conn) > 0 {
		c = conn[0]
	} else {
		c = nil
	}
	return &Player{ID: id, Update: make(chan bool), Conn: c, Name: name}
}

type Players map[string]*Player
