package werewolves

// Events names
const (
	EventTypeKill         = "kill"
	EventTypeStateChanged = "state_changed"
	EventTypeShowRole     = "show_role"

	EventTypeRevealRole = "reveal_role"

	EventTypeSavedFromDeath = "saved_from_death"
)

// Event represents an event that we put in Player.Update
// if u want to broadcast a kill event to
// all players in room then Example:
//
//	ev := &Event{EventTypeStateChanged, &StateChangedEvent{state}}
//	Room.BroadcastEvent(ev)
type Event struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data"`
}

func NewEvent(eventType string, data interface{}) Event {
	return Event{
		EventType: eventType,
		Data:      data,
	}
}

type (
	StateChangedEvent struct {
		State string `json:"state"`
	}
)
type EventRevealRole struct {
	Role     string `json:"role" mapstructure:"role"`
	PlayerID string `json:"player_id" mapstructure:"player_id"`
}

// TargetedEvent represents all events that require only a target id
type TargetedEvent struct {
	TargetID string `json:"player_id"`
}

type ConstableShootEvent struct {
	TargetID string `json:"target_id"`
	PlayerID string `json:"player_id"`
}

// FromEvent represents all event that require From player id and target id
type FromEvent struct {
	// player id who made a action
	FromID string `json:"from_id"`
	// target
	TargetID string `json:"target_id"`
}

type VoteEvent struct {
	// player id who made a action
	FromID string `json:"from_id"`
	// target
	TargetID string `json:"target_id"`
	// count of votes
	Votes uint8 `json:"votes"`
}
