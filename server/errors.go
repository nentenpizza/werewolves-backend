package server

var (
	AlreadyInRoomErr      = &ServerError{"you_in_room_error", ""}
	RoomNotFoundErr       = &ServerError{"room_not_found_error", ""}
	GameAlreadyStartedErr = &ServerError{"game_started_error", ""}
	PlayerNotFoundErr     = &ServerError{"player_not_found_error", ""}
	RoomStartErr          = &ServerError{"cannot_start_room_error", ""}
	NotAllowedErr         = &ServerError{"not_allowed_error", ""}
)

type ServerError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (se *ServerError) Error() string {
	return se.Type + " " + se.Message
}
