package websocket

var (
	RoomNotFoundErr       = &ServerError{Type: "room_not_found_error"}
	//GameAlreadyStartedErr = &ServerError{"game_started_error"}
	PlayerNotFoundErr     = &ServerError{Type: "player_not_found_error"}
	RoomStartErr          = &ServerError{Type: "cannot_start_room_error"}
	NotAllowedErr         = &ServerError{Type: "not_allowed_error"}

	NotInRoomRoom = &ServerError{Type: "you_not_in_room"}
	JoinRoomErr 		  =	&ServerError{Type: "failed_to_join_room"}
	AlreadyInRoomErr = &ServerError{Type: "already_in_room"}
	RoomNotExistsErr = &ServerError{Type: "room_not_exists"}
	RoomStartedErr = &ServerError{Type: "room_already_started"}
)

type ServerError struct {
	Type    string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (se *ServerError) Error() string {
	return se.Type
}