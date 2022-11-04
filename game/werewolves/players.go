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

	Groups []string `json:"groups"`

	sync.Mutex `json:"-"`
}

func (p *Player) Vote(pID string, votes uint8) Action {
	p.Lock()
	defer p.Unlock()
	return NewAction(
		VoteAction,

		func(r *Room) error {
			if r.Dead[p.ID] {
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

			if r.State == Discuss {
				return errors.New("game: can not vote in state discuss")
			}
			if r.State == Night {
				if p.Role == "Werewolf" || p.Role == "AlphaWerewolf" {
				} else {
					return errors.New("game: can not vote in night as long as you aren't werewolf")
				}
			}
			p.Voted = true
			r.votesCount++
			r.Votes[pID] += votes
			if r.votesCount == len(r.Players) {
				go r.forceNextState()
			}
			return nil
		},

		NewEvent(VoteAction, VoteEvent{p.ID, pID, votes}),
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
