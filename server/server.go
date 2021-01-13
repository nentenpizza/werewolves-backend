package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nentenpizza/werewolves/werewolves"
)

type Response map[string]interface{}

// Server represents a game server which talks with game
// PlayersRoom is map[PlayerID]RoomID
type Server struct {
	Rooms       map[string]*werewolves.Room
	PlayersRoom map[string]string
}

func (s *Server) HandleEvent(event *Event, conn *websocket.Conn) {
	switch event.Type {
	case EventTypeCreateRoom:
		ev := event.Data.(EventCreateRoom)
		s.handleCreateRoom(ev, conn)
	case EventTypeJoinRoom:
		ev := event.Data.(EventJoinRoom)
		s.handleJoinRoom(ev, conn)
	case EventTypeLeaveRoom:

	case EventTypeStartGame:
	}
}

func (s *Server) handleCreateRoom(event EventCreateRoom, conn *websocket.Conn) {
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, player.ID)
	s.Rooms[room.ID] = room
	s.PlayersRoom[player.ID] = room.ID
	go func() {
		conn.WriteJSON(player)
		for {
			select {
			case value := <-player.Update:
				if value == true {
					conn.WriteJSON(player)
				} else {
					return
				}
			}
		}
	}()
}

func (s *Server) handleJoinRoom(event EventJoinRoom, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(ServerError{RoomNotFound, ""})
		return
	}
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	err := room.AddPlayer(player)
	if err != nil {
		conn.WriteJSON(ServerError{GameAlreadyStarted, ""})
		return
	}
}

func (s *Server) handleLeaveRoom(event EventLeaveRoom, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(ServerError{RoomNotFound, ""})
		return
	}
	player, ok := room.Players[event.PlayerID]
	if !ok {
		conn.WriteJSON(ServerError{PlayerNotFound, ""})
		return
	}
	room.RemovePlayer(player.ID)
}

func (s *Server) handleStartGame(event EventStartGame, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(ServerError{RoomNotFound, ""})
		return
	}
	if event.PlayerID != room.Owner {
		conn.WriteJSON(ServerError{NotAllowed, ""})
		return
	}
	err := room.Run()
	if err != nil {
		conn.WriteJSON(ServerError{RoomStartErr, err.Error()})
		return
	}
}
