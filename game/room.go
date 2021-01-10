package game

import (
	"sync"
	"time"
)

// Werewolves does not have discuss time
// instead they voting just in night
const (
	DayVoting = "voting"
	Night = "night"
	Discuss = "discuss"
)

// Room represents a game room
// State: voting, discuss, night etc...
// all game logic here
type Room struct {
	Players Players
	state string
	ticker *time.Ticker
	done chan bool
	started bool
	sync.Mutex
}

func (r *Room) init(){
	r.state = Discuss
	r.ticker = time.NewTicker(10 * time.Second)
}

// IsDone returns is game ended or not
func (r *Room) IsDone() bool{
	_, ok := <-r.done
	return ok
}

// Run runs a main cycle of room (as a goroutine)
func (r *Room) Run(){
	if !r.started {
		r.init()
		go r.runCycle()
		r.started = true
	}
}

func (r *Room)  runCycle(){
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
func (r *Room) nextState(){
	r.Lock()
	defer r.Unlock()
	var nextState string
	switch r.state  {
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


// AddPlayer adds player to Room.Players
func (r *Room) AddPlayer(p Player){
	r.Lock()
	defer r.Unlock()
}