package werewolves

import (
	"errors"
	"sync"
)

type Player struct {
	Role string `json:"role"`

	// Character is a in-game character, hp, skills and other its a
	// character responsability
	Character Character `json:"character"`

	// Voted means player voted or not in current phase
	Voted bool `json:"voted"`

	// Unique Identifier of player
	ID string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	// Here we put an in-game updates
	Update chan Event `json:"-"`

	Room *Room `json:"-"`

	sync.Mutex `json:"-"`
}

func (p *Player) Vote(pID string) Action {
	p.Lock()
	defer p.Unlock()
	return NewAction(
		VoteAction,

		func(r *Room) error {
			if p.Voted == true {
				return errors.New("game: player already voted")
			}
			_, ok := r.Players[pID]
			if !ok {
				return errors.New("game: player is not in room")
			}

			if r.State == DayVoting {
				r.Votes[pID]++
				p.Voted = true
			} else if r.State == Night {
				if p.Role == "Werewolf" || p.Role == "AlphaWerewolf" { // мне насрать
					r.Votes[pID]++
					p.Voted = true
				} else {
					return errors.New("game: can not vote in night as long as you not werewolf")
				}
			} else {
				return errors.New("game: can not vote in state discuss")
			}

			return nil
		},

		NewEvent(VoteAction, FromEvent{p.ID, pID}),
	)
}

func (p *Player) Kill() {
	p.Character.SetHP(0)
	p.Room.Dead[p.ID] = true
}

func NewPlayer(id string, name string) *Player {
	return &Player{ID: id, Update: make(chan Event), Name: name}
}

type Players map[string]*Player
