package werewolves

import (
	"errors"
	"log"
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
	ID string `json:"id"`

	Name string `json:"name"`

	// Here we put an in-game updates
	Update chan interface{} `json:"-"`
	Room *Room `json:"room"`



	sync.Mutex `json:"-"`
}

func (p *Player) Vote(pID string) Action {
	p.Lock()
	defer p.Unlock()
	return NewAction(
		VoteAction,

		func(r *Room) error {
			if r.Dead[p.ID]{
				return errors.New("game: you are dead")
			}
			if r.Dead[pID] {
				return errors.New("game: target already dead")
			}
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
	p.Lock()
	defer p.Unlock()
	p.Character.SetHP(0)
	p.Room.Dead[p.ID] = true
	e := NewEvent(ExecutionAction, TargetedEvent{p.ID})
	p.Room.BroadcastEvent(e)
	log.Println(p.Name, "killed")
}

func NewPlayer(id string, name string) *Player {
	return &Player{ID: id, Update: make(chan interface{}), Name: name}
}

type Players map[string]*Player
