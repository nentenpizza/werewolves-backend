package werewolves

import (
	"errors"
	"sync"
)

type Player struct {
	Role       string      `json:"role"`
	Character  Character   `json:"character"`
	Voted      bool        `json:"voted"`
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Update     chan []byte `json:"-"`
	Room       *Room       `json:"room"`
	sync.Mutex `json:"-"`
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

func NewPlayer(id string, name string) *Player {
	return &Player{ID: id, Update: make(chan []byte), Name: name}
}

type Players map[string]*Player
