package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nentenpizza/werewolves/werewolves"
)

// Server represents a game server which talks with game
// PlayersRoom is map[PlayerID]RoomID
type Server struct {
	Rooms       map[string]*werewolves.Room
	PlayersRoom map[string]string
}

func New() *Server {
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
		s.HandleEvent(&ev, conn)
	}
}

func (s *Server) HandleEvent(event *Event, conn *websocket.Conn) {
	switch event.Type {
	case EventTypeCreateRoom:
		var ev EventCreateRoom
		b, err := json.Marshal(event.Data)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(b, &ev)
		if err != nil {
			log.Println(err)
			return
		}
		s.handleCreateRoom(&ev, conn)
	case EventTypeJoinRoom:
		ev := event.Data.(EventJoinRoom)
		s.handleJoinRoom(&ev, conn)
	case EventTypeLeaveRoom:

	case EventTypeStartGame:
	}
}

func (s *Server) handleCreateRoom(event *EventCreateRoom, conn *websocket.Conn) {
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, player.ID)
	s.Rooms[room.ID] = room
	s.PlayersRoom[player.ID] = room.ID
	err := conn.WriteJSON(player)
	if err != nil {
		log.Println(err)
		conn.WriteJSON(RoomNotFoundErr)
	}
	go func() {
		for {
			select {
			case value := <-player.Update:
				if value == true {
					err := conn.WriteJSON(player)
					if err != nil {
						log.Println(err)
					}
				} else {
					return
				}
			}
		}
	}()
}

func (s *Server) handleJoinRoom(event *EventJoinRoom, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(RoomNotFoundErr)
		return
	}
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, event.PlayerName)
	err := room.AddPlayer(player)
	if err != nil {
		err := conn.WriteJSON(GameAlreadyStartedErr)
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func (s *Server) handleLeaveRoom(event *EventLeaveRoom, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(RoomNotFoundErr)
		return
	}
	player, ok := room.Players[event.PlayerID]
	if !ok {
		conn.WriteJSON(PlayerNotFoundErr)
		return
	}
	room.RemovePlayer(player.ID)
}

func (s *Server) handleStartGame(event *EventStartGame, conn *websocket.Conn) {
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		conn.WriteJSON(RoomNotFoundErr)
		return
	}
	if event.PlayerID != room.Owner {
		conn.WriteJSON(NotAllowedErr)
		return
	}
	err := room.Run()
	if err != nil {
		conn.WriteJSON(RoomStartErr)
		return
	}
}
