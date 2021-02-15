package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/nentenpizza/werewolves/werewolves"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

// Server represents a game server which talks with game
// PlayersRoom is map[PlayerID]RoomID
type Server struct {
	Rooms       map[string]*werewolves.Room
	PlayersRoom map[string]string
}

func NewServer() *Server {
	return &Server{
		Rooms:       make(map[string]*werewolves.Room),
		PlayersRoom: make(map[string]string),
	}
}

func (s *Server) WsReader(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(string(msg))
		s.HandleEvent(&ev, conn) // handle error pls
	}
}

func (s *Server) HandleEvent(event *Event, conn *websocket.Conn) error {
	switch event.Type {
	case EventTypeCreateRoom:
		ev := EventCreateRoom{}
		mapstructure.Decode(event.Data, &ev)
		err := s.handleCreateRoom(&ev, conn)
		if err != nil {
			return err
		}

	case EventTypeJoinRoom:
		ev := EventJoinRoom{}
		mapstructure.Decode(event.Data, &ev)
		err := s.handleJoinRoom(&ev, conn)
		if err != nil {
			return err
		}
	case EventTypeLeaveRoom:
		ev := EventLeaveRoom{}
		mapstructure.Decode(event.Data, &ev)
		err := s.handleLeaveRoom(&ev, conn)
		if err != nil {
			return err
		}
	case EventTypeStartGame:
		ev := EventStartGame{}
		mapstructure.Decode(event.Data, &ev)
		err := s.handleStartGame(&ev, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) handleCreateRoom(event *EventCreateRoom, conn *websocket.Conn) error {
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, player.ID)
	s.Rooms[room.ID] = room
	s.PlayersRoom[player.ID] = room.ID
	err := conn.WriteJSON(player)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case value := <-player.Update:
				err := conn.WriteJSON(value)
				if err != nil {
					return
				}
			}
		}
	}()
	return nil
}

func (s *Server) handleJoinRoom(event *EventJoinRoom, conn *websocket.Conn) error {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		return RoomNotFoundErr
	}
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	err := room.AddPlayer(player)

	return err
}

func (s *Server) handleLeaveRoom(event *EventLeaveRoom, conn *websocket.Conn) error {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		return RoomNotFoundErr
	}
	player, ok := room.Players[event.PlayerID]
	if !ok {
		return PlayerNotFoundErr
	}
	room.RemovePlayer(player.ID)
	return nil
}

func (s *Server) handleStartGame(event *EventStartGame, conn *websocket.Conn) error {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		return RoomNotFoundErr

	}
	if event.PlayerID != room.Owner {
		return NotAllowedErr
	}
	err := room.Run()
	if err != nil {
		return RoomStartErr
	}
	return nil
}
