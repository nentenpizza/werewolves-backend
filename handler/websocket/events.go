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
	EventTypeAllRooms   = "all_rooms"

	EventTypeRoomCreated = "new_room"
	EventTypeRoomDeleted = "room_deleted"
	EventTypeEndGame     = "end_game"

	EventTypeDisconnected = "disconnected"
	EventTypeNotInGame    = "not_in_game"

	EventTypeUseSkill   = "use_skill"
	EventTypeRevealRole = "reveal_role"

	EventTypeSendMessage = "send_message"
	EventTypeSendEmote   = "send_emote"
	EventTypeFloodWait   = "flood_wait"
	EventTypeEmojiWait   = "emoji_wait"
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

type EventErr struct {
	Type  string       `json:"event_type" mapstructure:"event_type"`
	Data  interface{}  `json:"data" mapstructure:"data"`
	Error *ServerError `json:"error" mapstructure:"error"`
}

// Event for pre-game stuff
type (
	// EventCreateRoom represents event for creating room
	EventCreateRoom struct {
		RoomName string              `json:"room_name" mapstructure:"room_name"`
		Settings werewolves.Settings `json:"settings" mapstructure:"settings"`
	}

	EventRoomPlayer struct {
		RoomID   string `json:"room_id,omitempty" mapstructure:"room_id"`
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
	}

	EventAllRooms struct {
		Rooms *Rooms `json:"rooms" mapstructure:"rooms"`
	}

	EventNewRoomCreated struct {
		Room *werewolves.Room `json:"room" mapstructure:"room"`
	}

	EventRoomDeleted struct {
		RoomID string `json:"room_id" mapstructure:"room_id"`
	}
)

type EventRevealRole struct {
	Role     string `json:"role" mapstructure:"role"`
	PlayerID string `json:"player_id" mapstructure:"player_id"`
}

type EventEndGame struct {
	WonGroup  map[string]*werewolves.Player `json:"won" mapstructure:"won"`
	LoseGroup map[string]*werewolves.Player `json:"lose" mapstructure:"lose"`

	// XP is the amount of experience the player received for this game
	XP int `json:"XP" mapstructure:"XP"`
}

// Events for chat
type (
	MessageEvent struct {
		Text     string `json:"text" mapstructure:"text"`
		Username string `json:"username,omitempty" mapstructure:"username"`
	}

	EmoteEvent struct {
		FromID string `json:"from_id,omitempty" mapstructure:"from_id,omitempty"`
		Emote  string `json:"emote" mapstructure:"emote"`
	}
)

type EventFloodWait struct {
	Left int64 `json:"left"`
}

// Events for in-game stuff
type (

	// TargetedEvent used in all cases when you need only player_id and target_id
	TargetedEvent struct {
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
		TargetID string `json:"target_id,omitempty" mapstructure:"target_id"`
	}

	EventPlayerID struct {
		PlayerID string `json:"player_id,omitempty" mapstructure:"player_id"`
	}
)
