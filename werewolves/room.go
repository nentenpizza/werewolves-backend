package werewolves

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sort"
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
	Prepare = "prepare"
)

// Duration of each phase
const (
	PhaseLength = 5 * time.Second
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
// Votes is map[PlayerID]Count of votes
//
// Broadcast is channel for sending a signal to all players in the room
// to make it clear that the state needs to be updated
// Owner is playerID to creator of room
//
type Room struct {
	Owner     string  `json:"owner"`
	Players   Players `json:"-"`
	Users map[string]bool `json:"players"`
	State     string  `json:"state"`
	ticker    *time.Ticker
	notifyDone chan bool
	done      chan bool
	started   bool
	Dead      map[string]bool   `json:"dead"`
	broadcast chan interface{}
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	OpenRoles map[string]string `json:"open_roles"`
	Votes     map[string]uint8  `json:"votes"`
	Settings  Settings          `json:"settings"`
	sync.Mutex
}


// NewRoom constructor for Room
// Pass empty Settings for defaults
func NewRoom(id string, name string, players Players, settings Settings, ownerID string) *Room {
	r := &Room{
		Players:   players,
		Users: make(map[string]bool),
		Dead:      make(map[string]bool),
		broadcast: make(chan interface{}),
		Settings:  settings,
		OpenRoles: make(map[string]string),
		Votes:     make(map[string]uint8),
		ID:        id,
		Name:      name,
		Owner:     ownerID,
		State: Prepare,
	}
	return r
}

func (r *Room) Done() chan bool {
	return r.notifyDone
}

func (r *Room) init() error {
	for _, p := range r.Players {
		p.Room = r
	}
	err := r.defineRoles()
	if err != nil {
		return err
	}
	r.State = Discuss
	r.done = make(chan bool, 1)
	r.ticker = time.NewTicker(PhaseLength)
	for k := range r.Players {
		r.Votes[k] = 0
	}
	return nil
}

func (r *Room) Started() bool {
	return r.started
}

// Run define roles and runs a main cycle of room (as a goroutine)
func (r *Room) Start() error {
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
		go r.runCycle()
		r.started = true
		r.nextState()
	}
	return nil
}

func (r *Room) runCycle() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			r.notifyDone <- true
			return

		case <-r.ticker.C:
			r.nextState()

		}
	}
}

func (r *Room) endVotePhase() {
	type kv struct {
		Key   string
		Value uint8
	}
	var sortedVotes []kv
	for k, v := range r.Votes {
		sortedVotes = append(sortedVotes, kv{k, v})
	}
	sort.Slice(sortedVotes, func(i, j int) bool {
		return sortedVotes[i].Value > sortedVotes[j].Value
	})

	if sortedVotes[0].Value  == sortedVotes[1].Value  {
		return
	} else if sortedVotes[0].Value < uint8(len(r.Players)/2) && len(r.Players) > 2{
		return
	} else {
		p, ok := r.Players[sortedVotes[0].Key]
		if ok {
			p.Kill()
		}
	}

	r.resetVotes()
}

func (r *Room) resetVotes() {
	for k := range r.Votes {
		r.Votes[k] = 0
	}
	for _, p := range r.Players {
		p.Voted = false
	}
}

// Changes state to next value in game loop
func (r *Room) nextState() {
	r.Lock()
	defer r.Unlock()
	switch r.State {
	case Discuss:
		r.State = DayVoting
	case DayVoting:
		r.State = Night
		r.endVotePhase()
	case Night:
		r.State = Discuss
		r.endVotePhase()
		r.resetProtection()
	default:
		break
	}
	if len(r.Players) < MinPlayers{
		r.done <- true
	}
	log.Println(r.State)
	ev := Event{EventTypeStateChanged, r}
	r.BroadcastEvent(ev)
}

// Perform validates Action and performs it
func (r *Room) Perform(action Action) error {
	r.Lock()
	defer r.Unlock()
	var ok bool
	actions := allowedActions[r.State]
	for _, v := range actions {
		if v == action.Name {
			ok = true
		}
	}
	if !ok {
		return errors.New("game: action not allowed")
	}
	err := action.do(r)
	if err == nil {
		r.BroadcastEvent(action.Event)
	}
	log.Println(action.Name, "performed")
	r.commitDead()

	return err
}

func (r *Room) commitDead() {
	for _, p := range r.Players {
		if p.Character.IsDead() {
			r.Dead[p.ID] = true
			if r.Settings.OpenRolesOnDeath {
				r.OpenRoles[p.ID] = p.Role
			}
		}
	}
}

// AddPlayer adds player to Room.Players
func (r *Room) AddPlayer(p *Player) error {
	r.Lock()
	defer r.Unlock()
	if !r.started {
		r.Players[p.ID] = p
		r.Users[p.Name] = true
		p.Room = r
	} else {
		return errors.New("game: can't add player to started room")
	}
	return nil
}

// RemovePlayer Removes player from room
// if game started, then player will be killed
func (r *Room) RemovePlayer(playerID string) error {
	r.Lock()
	defer r.Unlock()
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
		delete(r.Players, p.ID)
		kill := TargetedEvent{TargetID: p.ID}
		r.BroadcastEvent(Event{EventType: ExecutionAction, Data: kill})
	}
	return nil
}

func (r *Room) Resurrect(playerID string) error {
	p, ok := r.Players[playerID]
	if !ok {
		return fmt.Errorf(
			"game: player with id: %s is not in %s room, room_id: %s", playerID, r.Name, r.ID,
		)
	}
	p.Character.SetHP(1)
	delete(r.Dead, playerID)
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
		role := roles[i](v.ID)
		v.Character = role

		v.Role = reflect.TypeOf(role).Elem().Name()
		v.Update <- Event{EventType: EventTypeShowRole, Data: struct {
			Name string `json:"name"`
			Character Character `json:"character"`
		}{Name: v.Role, Character: v.Character}}
		i++
	}
	return nil
}


// resets all doctor and other protection
func (r *Room) resetProtection() {
	for _, p := range r.Players {
		if p.Character.HP() > 1 {
			p.Character.SetHP(1)
		}
	}
}


func (r *Room) BroadcastEvent(e interface{}) {
	for _, p := range r.Players {
		p.Update <- e
	}
}

func (r *Room) Broadcast() chan interface{} {
	return r.broadcast
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
