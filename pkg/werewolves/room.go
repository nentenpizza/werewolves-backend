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
// instead they voting only on night
const (
	DayVoting = "voting"
	Night     = "night"
	Discuss   = "discuss"
	Prepare   = "prepare"
)

// Duration of each phase
var (
	PhaseLength = 5 * time.Second
)

// Capacity of room
const (
	MaxPlayers = 12
	MinPlayers = 2 // 6
)

type RoomResult struct {
	WonGroup string `json:"won_group,omitempty" mapstructure:"won_group,omitempty"`
}

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
type Room struct {
	Owner       string          `json:"owner"`
	Players     Players         `json:"-"`
	Users       map[string]bool `json:"players"`
	State       string          `json:"state"`
	ticker      *time.Ticker
	NotifyDone  chan *RoomResult `json:"-"`
	done        chan *RoomResult
	started     bool
	Dead        map[string]bool `json:"dead"`
	broadcast   chan interface{}
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	OpenRoles   map[string]string `json:"open_roles"`
	Votes       map[string]uint8  `json:"votes"`
	votesCount  int
	Settings    Settings                      `json:"settings"`
	Groups      map[string]map[string]*Player `json:"-"`
	AliveGroups map[string]map[string]*Player `json:"-"`
	sync.Mutex
}

// NewRoom constructor for Room
// Pass empty Settings for defaults
func NewRoom(id string, name string, players Players, settings Settings, ownerID string) *Room {
	r := &Room{
		Players:     players,
		Users:       make(map[string]bool),
		Dead:        make(map[string]bool),
		broadcast:   make(chan interface{}),
		Settings:    settings,
		OpenRoles:   make(map[string]string),
		Votes:       make(map[string]uint8),
		ID:          id,
		Name:        name,
		Owner:       ownerID,
		State:       Prepare,
		Groups:      make(map[string]map[string]*Player),
		AliveGroups: make(map[string]map[string]*Player),
		done:        make(chan *RoomResult, 1),
		NotifyDone:  make(chan *RoomResult, 1),
	}
	return r
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
	r.ticker = time.NewTicker(PhaseLength)
	for k := range r.Players {
		r.Votes[k] = 0
	}
	return nil
}

func (r *Room) Started() bool {
	return r.started
}

// Start defines roles and runs main cycle of room (as a goroutine)
func (r *Room) Start() error {
	if !r.started {
		if len(r.Players) < MinPlayers {
			return NotEnoughPlayers
		}
		if len(r.Players) > MaxPlayers {
			return errors.New("too much players in room")
		}
		err := r.init()
		if err != nil {
			return err
		}
		r.revealTeams()
		go r.runCycle()
		r.started = true
		r.nextState()
	}
	return nil
}

func (r *Room) runCycle() {
	for {
		select {
		case e := <-r.done:
			r.ticker.Stop()
			r.NotifyDone <- e
			return

		case <-r.ticker.C:
			r.Lock()
			r.nextState()
			r.Unlock()
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

	if sortedVotes[0].Value == sortedVotes[1].Value {
		return
	} else if sortedVotes[0].Value < uint8(len(r.Players)/2) && len(r.Players) > 2 {
		return
	} else {
		p, ok := r.Players[sortedVotes[0].Key]
		if ok {
			r.doKillPlayer(p)
		}
	}

}

func (r *Room) resetVotes() {
	for k := range r.Votes {
		r.Votes[k] = 0
	}
	for _, p := range r.Players {
		p.Voted = false
	}
	r.votesCount = 0
}

// Changes state to next value in game loop
func (r *Room) nextState() {
	switch r.State {
	case Discuss:
		r.State = DayVoting
	case DayVoting:
		r.State = Night
		r.endVotePhase()
		r.resetVotes()
	case Night:
		r.State = Discuss
		r.endVotePhase()
		r.resetVotes()
		r.resetProtection()
	default:
		break
	}
	if len(r.Players) < MinPlayers {
		r.done <- nil
	}
	log.Println(r.State)
	ev := Event{EventTypeStateChanged, r}
	r.BroadcastEvent(ev)
}

func (r *Room) forceNextState() {
	r.ticker.Reset(PhaseLength)
	r.nextState()
}

// Perform validates Action and performs it
// then sends Action.Event to all players or players in specific groups
func (r *Room) Perform(action Action, groups ...string) error {
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
		if len(groups) > 0 {
			for _, v := range groups {
				r.doBroadcastTo(v, action.Event)
			}
		} else {
			r.BroadcastEvent(action.Event)
		}
	}
	log.Println(action.Name, "performed")
	r.commitDead()

	return err
}

func (r *Room) commitDead() {
	for _, p := range r.Players {
		if p.Character.IsDead() {
			r.doKillPlayer(p)
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
		r.doKillPlayer(p)
		delete(r.Players, p.ID)
		kill := TargetedEvent{TargetID: p.ID}
		r.BroadcastEvent(Event{EventType: ExecutionAction, Data: kill})
	}
	return nil
}

func (r *Room) Resurrect(playerID string) error {
	r.Lock()
	defer r.Unlock()
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
	roles, err := genRolesList(len(r.Players))
	if err != nil {
		return err
	}
	var i int
	for _, v := range r.Players {
		role := roles[i](v.ID)
		v.Character = role
		v.Role = reflect.TypeOf(role).Elem().Name()
		if v.Role == "Werewolf" || v.Role == "AlphaWerewolf" {
			r.doJoinGroup("wolves", v)
			v.Groups = append(v.Groups, "wolves")
		} else {
			r.doJoinGroup("peaceful", v)
			v.Groups = append(v.Groups, "peaceful")
		}
		v.Update <- Event{
			EventType: EventTypeShowRole, Data: struct {
				Name      string    `json:"name"`
				Character Character `json:"character"`
			}{Name: v.Role, Character: v.Character}}
		i++
	}
	return nil
}

func (r *Room) revealTeams() {
	wolves, ok := r.Groups["wolves"]
	if ok {
		for _, p := range wolves {
			r.doBroadcastTo("wolves", &Event{EventTypeRevealRole, &EventRevealRole{Role: p.Role, PlayerID: p.ID}})
		}
	}
}

// InGroup returns true if player with playerID in group with groupName
func (r *Room) InGroup(groupName string, playerID string) bool {
	group, ok := r.Groups[groupName]
	if ok {
		for _, p := range group {
			if p.ID == playerID {
				return true
			}
		}
	}
	return false
}

func (r *Room) JoinGroup(groupName string, p *Player) {
	r.Lock()
	defer r.Unlock()
	r.doJoinGroup(groupName, p)

}

func (r *Room) doJoinGroup(groupName string, p *Player) {
	group, ok := r.Groups[groupName]
	if ok {
		group[p.ID] = p
	} else {
		r.Groups[groupName] = map[string]*Player{p.ID: p}
	}
	if groupName != "dead" {
		r.doJoinAliveGroup(groupName, p)
	}
}

func (r *Room) doJoinAliveGroup(groupName string, p *Player) {
	group, ok := r.AliveGroups[groupName]
	if ok {
		group[p.ID] = p
	} else {
		r.AliveGroups[groupName] = map[string]*Player{p.ID: p}
	}
}

func (r *Room) removeFromAliveGroup(groupName string, p *Player) {
	group, ok := r.AliveGroups[groupName]
	if ok {
		delete(group, p.ID)
	}
}

func (r *Room) BroadcastTo(groupName string, i interface{}) {
	r.Lock()
	defer r.Unlock()
	r.doBroadcastTo(groupName, i)
}

func (r *Room) doBroadcastTo(groupName string, i interface{}) {
	group, ok := r.Groups[groupName]
	if ok {
		for _, p := range group {
			p.Update <- i
		}
	}
}

func (r *Room) KillPlayer(player *Player) {
	r.Lock()
	defer r.Unlock()
	r.doKillPlayer(player)
}

func (r *Room) doKillPlayer(player *Player) {
	if player == nil {
		return
	}
	if player.Character.HP() <= 1 {
		wolves, ok := r.AliveGroups["wolves"]
		peaceful, okk := r.AliveGroups["peaceful"]
		player.Kill()
		for _, v := range player.Groups {
			r.removeFromAliveGroup(v, player)
		}
		r.doJoinGroup("dead", player)
		if ok && okk {
			if len(wolves) >= len(peaceful) {
				r.done <- &RoomResult{WonGroup: "wolves"}
			}
			if len(wolves) <= 0 {
				r.done <- &RoomResult{WonGroup: "peaceful"}
			}
		}
	} else {
		r.BroadcastEvent(&Event{EventTypeSavedFromDeath, TargetedEvent{TargetID: player.ID}})
		player.Character.SetHP(1)
	}
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
