package game

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// States
// Werewolves does not have discuss time
// instead they voting just in night
const (
	DayVoting = "voting"
	Night     = "night"
	Discuss   = "discuss"
)

const (
	MaxPlayers = 10
	MinPlayers = 2 // 6
)

// Room represents a game room
// State: voting, discuss, night etc...
// all game logic here
type Room struct {
	Players   Players `json:"players"`
	state     string
	ticker    *time.Ticker
	done      chan bool
	started   bool
	Dead      Players `json:"dead"`
	Broadcast chan bool
	sync.Mutex
}

func NewRoom(players Players) *Room {
	return &Room{Players: players, Dead: Players{}, Broadcast: make(chan bool)}
}

func (r *Room) init() {
	for _, p := range r.Players {
		p.Room = r
	}
	r.defineRoles()
	r.state = Discuss
	r.done = make(chan bool, 1)
	r.ticker = time.NewTicker(10 * time.Second)
}

// IsDone returns is game ended or not
func (r *Room) IsDone() bool {
	_, ok := <-r.done
	return ok
}

// Run define roles
// and runs a main cycle of room (as a goroutine)
func (r *Room) Run() error {
	if !r.started {
		if len(r.Players) < MinPlayers {
			return errors.New("not enough players to start")
		}
		if len(r.Players) > MaxPlayers {
			return errors.New("too much players in room")
		}
		r.init()
		go func() {
			select {
			case <-r.Broadcast:
				for _, p := range r.Players {
					p.Update <- true
				}
			}
		}()
		go r.runCycle()
		r.started = true
	}
	return nil
}

func (r *Room) runCycle() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			break

		case <-r.ticker.C:
			r.nextState()

		}
	}
}

// Changes state to next value in game loop
func (r *Room) nextState() {
	r.Lock()
	defer r.Unlock()
	switch r.state {
	case Discuss:
		r.state = DayVoting
	case DayVoting:
		r.state = Night
	case Night:
		r.state = Discuss
		r.resetProtection()
	default:
		break
	}
}

// Perform validates Action and performs it
func (r *Room) Perform(action Action) error {
	r.Lock()
	defer r.Unlock()
	var ok bool
	actions := allowedActions[r.state]
	for _, v := range actions {
		if v == action.Name {
			ok = true
		}
	}
	if !ok {
		return errors.New("game: action not allowed")
	}
	action.do(r)
	for _, p := range r.Players {
		if p.Character.Dead() {
			r.Dead[p.ID] = p

		}
		p.Update <- true
	}
	return nil
}

// AddPlayer adds player to Room.Players
func (r *Room) AddPlayer(p *Player) {
	r.Lock()
	defer r.Unlock()
	r.Players[p.ID] = p
}

// defineRoles defines roles for Room.Players
func (r *Room) defineRoles() {
	roles, err := rolesList(len(r.Players))
	if err != nil {
		panic(err)
	}
	var i int
	for _, v := range r.Players {
		v.Character = roles[i]()
		i++
	}

}

// resets all doctor etc... protection
func (r *Room) resetProtection() {
	for _, p := range r.Players {
		p.Character.SetHP(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
