package game

import (
	"errors"
	"fmt"
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

// Capacity of room
const (
	MaxPlayers = 10
	MinPlayers = 2 // 6
)

// Settings for Room, not required
type Settings struct {
	OpenRolesOnDeath bool `json:"open_roles_on_death"`
}

// Room represents a game room
// State: voting, discuss, night etc...
// all game logic here
//
// OpenRoles is map[PlayerID]RoleName
// Dead is map[PlayerID]bool
//
// Broadcast is channel for sending a signal to all players in the room
// to make it clear that the state needs to be updated
type Room struct {
	Players   Players `json:"-"`
	state     string
	ticker    *time.Ticker
	done      chan bool
	started   bool
	Dead      map[string]bool   `json:"dead"`
	Broadcast chan bool         `json:"-"`
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	OpenRoles map[string]string `json:"open_roles"`
	Settings  Settings          `json:"settings"`
	sync.Mutex
}

// NewRoom constructor for Room
// Pass empty Settings for defaults
func NewRoom(id string, name string, players Players, settings Settings) *Room {
	return &Room{
		Players:   players,
		Dead:      make(map[string]bool),
		Broadcast: make(chan bool),
		Settings:  settings,
		OpenRoles: make(map[string]string)}
}

func (r *Room) init() error {
	for _, p := range r.Players {
		p.Room = r
	}
	err := r.defineRoles()
	if err != nil {
		return err
	}
	r.state = Discuss
	r.done = make(chan bool, 1)
	r.ticker = time.NewTicker(10 * time.Second)
	return nil
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
		err := r.init()
		if err != nil {
			return err
		}
		go r.runBroadcaster()
		go r.runCycle()
		r.started = true
	}
	return nil
}

func (r *Room) runBroadcaster() {
	select {
	case <-r.Broadcast:
		for _, p := range r.Players {
			p.Update <- true
		}
	}
}

func (r *Room) runCycle() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			return

		case <-r.ticker.C:
			r.nextState()
			r.refreshPlayers()

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
	r.appendDead()
	r.refreshPlayers()
	return nil
}

func (r *Room) appendDead() {
	for _, p := range r.Players {
		if p.Character.IsDead() {
			r.Dead[p.ID] = true
		}
	}
}

// AddPlayer adds player to Room.Players
func (r *Room) AddPlayer(p *Player) error {
	r.Lock()
	defer r.Unlock()
	if !r.started {
		r.Players[p.ID] = p
	} else {
		return errors.New("game: can't add player to started room")
	}
	return nil
}

// RemovePlayer Removes player from room
// if game started, then player will be killed
func (r *Room) RemovePlayer(playerID string) error {
	p, ok := r.Players[playerID]
	if !ok {
		return fmt.Errorf(
			"game: player with id: %s is not in %s room, room_id: %s", playerID, r.Name, r.ID,
		)
	}
	if !r.started {
		delete(r.Players, p.ID)
	} else {
		p.Kill()
	}
	r.refreshPlayers()
	return nil
}

// defineRoles defines roles for Room.Players
func (r *Room) defineRoles() error {
	roles, err := rolesList(len(r.Players))
	if err != nil {
		return err
	}
	var i int
	for _, v := range r.Players {
		v.Character = roles[i]()
		i++
	}
	return nil
}

// resets all doctor and other protection
func (r *Room) resetProtection() {
	for _, p := range r.Players {
		p.Character.SetHP(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (r *Room) refreshPlayers() {
	r.Broadcast <- true
}
