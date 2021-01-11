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
	MinPlayers = 6
)

// Room represents a game room
// State: voting, discuss, night etc...
// all game logic here
type Room struct {
	Players Players
	Room    *Room
	state   string
	ticker  *time.Ticker
	done    chan bool
	started bool
	sync.Mutex
}

func NewRoom(players Players) *Room {
	return &Room{Players: players}
}

func (r *Room) init() {
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
	var nextState string
	switch r.state {
	case Discuss:
		nextState = DayVoting
	case DayVoting:
		nextState = Night
	case Night:
		nextState = Discuss
	default:
		nextState = r.state
	}
	r.state = nextState
}

// Perform validates Action and performs it
func (r *Room) Perform(action Action) error {
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
		v.Role = roles[i]()
		i++
	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
