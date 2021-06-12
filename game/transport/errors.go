package transport

var (
	RoomNotFoundErr = &ServerError{Type: "RoomNotFoundErr"}
	//GameAlreadyStartedErr = &ServerError{"game_started_error"}
	PlayerNotFoundErr = &ServerError{Type: "PlayerNotFoundErr"}
	RoomStartErr      = &ServerError{Type: "RoomStartErr"}
	NotAllowedErr     = &ServerError{Type: "NotAllowedErr"}

	NotInRoomRoom    = &ServerError{Type: "NotInRoomRoom"}
	JoinRoomErr      = &ServerError{Type: "JoinRoomErr"}
	AlreadyInRoomErr = &ServerError{Type: "AlreadyInRoomErr"}
	RoomNotExistsErr = &ServerError{Type: "RoomNotExistsErr"}
	RoomStartedErr   = &ServerError{Type: "RoomStartedErr"}
)

type ServerError struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}

func (se *ServerError) Error() string {
	return se.Type
}
