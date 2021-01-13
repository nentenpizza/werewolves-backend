package server

import (
	"github.com/nentenpizza/werewolves/werewolves"
)

// Event Types for typical things
const (
	EventTypeCreateRoom = "create_room"
	EventTypeJoinRoom   = "join_room"
	EventTypeLeaveRoom  = "leave_room"
	EventTypeVote       = werewolves.VoteAction
)

// Event types for skills
const (
	EventTypeConstableShoot = werewolves.ConstableShootAction
	EventTypeDoctorHeal     = werewolves.DoctorHealAction
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Event for pre-game stuff
type (
	// EventCreateRoom represents event for creating room
	EventCreateRoom struct {
		Name     string              `json:"name"`
		ID       string              `json:"id"`
		Settings werewolves.Settings `json:"settings"`
	}

	// EventJoinRoom represents event for joining room
	EventJoinRoom struct {
		PlayerName string `json:"player_name"`
	}

	// EventLeaveRoom represents event for leaving room
	EventLeaveRoom struct {
		PlayerID string `json:"player_id"`
	}
)

// Events for in-game stuff
type (
	// used in all cases when you need only player_id and target_id
	TargetedEvent struct {
		PlayerID string `json:"player_id`
		TargetID string `json:"player_id`
	}
)
