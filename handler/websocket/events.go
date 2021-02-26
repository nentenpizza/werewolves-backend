package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
)

// Event Types for typical things
const (
	EventTypeCreateRoom = "create_room"
	EventTypeJoinRoom   = "join_room"
	EventTypeLeaveRoom  = "leave_room"
	EventTypeStartGame  = "start_room"
	EventTypeVote       = werewolves.VoteAction

	EventTypeSendMessage = "send_message"
)

// Event types for skills
const (
	EventTypeConstableShoot = werewolves.ConstableShootAction
	EventTypeDoctorHeal     = werewolves.DoctorHealAction
)


type Event struct {
	Type string      `json:"event_type" mapstructure:"event_type"`
	Data interface{} `json:"data" mapstructure:"data"`
}

// Event for pre-game stuff
type (
	// EventCreateRoom represents event for creating room
	EventCreateRoom struct {
		RoomName   string              `json:"room_name" mapstructure:"room_name"`
		Settings   werewolves.Settings `json:"settings" mapstructure:"settings"`
	}

	// EventLeaveRoom represents event for leaving room
	EventLeaveRoom struct {
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
		RoomID   string `json:"room_id,omitempty" mapstructure:"room_id"`
	}


	EventJoinRoom struct {
		RoomID string `json:"room_id" mapstructure:"room_id"`
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
	}

)

// Events for chat
type (
	MessageEvent struct {
		Text string `json:"text" mapstructure:"text"`
		Username string `json:"username,omitempty" mapstructure:"username"`
	}
)

// Events for in-game stuff
type (


	// TargetedEvent used in all cases when you need only player_id and target_id
	TargetedEvent struct {
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
		TargetID string `json:"target_id,omitempty" mapstructure:"target_id"`
	}
)

var EventsWithTypes = map[string]interface {}{
	EventTypeStartGame: Event{},
}