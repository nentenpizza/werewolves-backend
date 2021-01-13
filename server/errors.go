package server

const (
	AlreadyInRoom      = "you_in_room_error"
	RoomNotFound       = "room_not_found_error"
	GameAlreadyStarted = "game_started_error"
	PlayerNotFound     = "player_not_found_error"
	RoomStartErr       = "cannot_start_room_error"
	NotAllowed         = "not_allowed_error"
)

type ServerError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
