package handler

import (
	"encoding/json"
	"errors"
	j"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
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
	handler
	Rooms       map[string]*werewolves.Room
	PlayersRoom map[string]string
	Secret []byte
}

func (s Server) REGISTER(h handler, g *echo.Group)  {
	for i:=0;i<10;i++ {
		room := werewolves.NewRoom(uuid.New().String(), uuid.New().String(), werewolves.Players{}, werewolves.Settings{}, uuid.New().String())
		s.Rooms[room.ID] = room
	}
	s.handler = h
	g.GET("/allrooms", s.AllRooms)
}

func NewServer(secret []byte) *Server {
	return &Server{
		Rooms:       make(map[string]*werewolves.Room),
		PlayersRoom: make(map[string]string),
		Secret: secret,
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
		token, err := j.ParseWithClaims(ev.Token, &jwt.Claims{}, func(token *j.Token) (interface{}, error) {
			return s.Secret, nil
		})
		if err != nil {
			return
		}
		if !token.Valid {
			continue
		}
		jwtWithClaims := jwt.From(token)
		err = s.HandleEvent(&ev, conn, jwtWithClaims) // handle error pls
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) HandleEvent(event *Event, conn *websocket.Conn, token jwt.Claims) error {
	switch event.Type {
	case EventTypeCreateRoom:
		ev := EventCreateRoom{}
		err := mapstructure.Decode(event.Data, &ev)
		if err != nil {
			return err
		}
		err = s.handleCreateRoom(&ev, conn, token)
		if err != nil {
			return err
		}

	case EventTypeJoinRoom:
		ev := &EventJoinRoom{}
		b, err  := json.Marshal(event.Data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, ev)
		if err != nil {
			return err
		}
		err = s.handleJoinRoom(ev, conn, token)
		if err != nil {
			return err
		}
	case EventTypeLeaveRoom:
		ev := EventLeaveRoom{}
		err := mapstructure.Decode(event.Data, &ev)
		if err != nil {
			return err
		}
		err = s.handleLeaveRoom(&ev, conn)
		if err != nil {
			return err
		}
	case EventTypeStartGame:
		ev := EventStartGame{}
		err := mapstructure.Decode(event.Data, &ev)
		if err != nil {
			return err
		}
		err = s.handleStartGame(&ev, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) handleCreateRoom(event *EventCreateRoom, conn *websocket.Conn, token jwt.Claims) error {
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, token.Username)
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


func (s *Server) handleJoinRoom(event *EventJoinRoom, conn *websocket.Conn, token jwt.Claims) error {
	_, ok := s.PlayersRoom[token.Username]
	if ok{
		return errors.New("player already in room")
	}
	room, ok := s.Rooms[event.RoomID]
	if !ok {
		return RoomNotFoundErr
	}
	s.PlayersRoom[token.Username] = room.Name
	id := uuid.New().String()
	player := werewolves.NewPlayer(id, token.Username)
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
	err := room.RemovePlayer(player.ID)
	return err
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
