package handler

import (
	"encoding/json"
	j "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nentenpizza/werewolves/werewolves"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type Player struct {
	sync.Mutex
	conn *websocket.Conn
	*werewolves.Player
	room *werewolves.Room
}

func (p *Player) Conn() *websocket.Conn {
	p.Lock()
	defer p.Unlock()
	return p.conn
}

func (p *Player) SetRoom(r *werewolves.Room){
	p.Lock()
	defer p.Unlock()
	p.room = r
}
func (p *Player) SetChar(plr *werewolves.Player){
	p.Lock()
	defer p.Unlock()
	p.Player = plr
}

func (p *Player) Room()*werewolves.Room{
	p.Lock()
	defer p.Unlock()
	return p.room
}
func (p *Player) Char() *werewolves.Player{
	p.Lock()
	defer p.Unlock()
	return p.Player
}

// Server represents a game server which talks with game
// PlayersRoom is map[PlayerID]RoomID
type Server struct {
	handler
	Rooms       map[string]*werewolves.Room
	PlayersRoom map[string]string
	Players map[string]*Player
	Secret []byte
}

func (s Server) REGISTER(h handler, g *echo.Group)  {
	for i:=0;i<5;i++ {
		room := werewolves.NewRoom(strconv.Itoa(rand.Intn(100)), strconv.Itoa(rand.Intn(100)), werewolves.Players{}, werewolves.Settings{}, strconv.Itoa(rand.Intn(100)))
		s.Rooms[room.ID] = room
	}
	s.handler = h
	g.GET("/allrooms", s.AllRooms)
}

func NewServer(secret []byte) *Server {
	return &Server{
		Rooms:       make(map[string]*werewolves.Room),
		PlayersRoom: make(map[string]string),
		Players: make(map[string]*Player),
		Secret: secret,
	}
}

func (s *Server) WsReader(conn *websocket.Conn) {
	var username string
	defer conn.Close()
	defer delete(s.Players, username)
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
		if username == ""{
			username = jwtWithClaims.Username
			s.Players[jwtWithClaims.Username] = &Player{conn: conn}
		}
		err = s.HandleEvent(&ev, conn, jwtWithClaims) // handle error pls
		if err != nil {
			log.Println(err)
			conn.WriteJSON(err)
		}
	}
}

func (s *Server) HandleEvent(event *Event, conn *websocket.Conn, token jwt.Claims) error {
	js, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	switch event.Type {
	case EventTypeCreateRoom:
		ev := EventCreateRoom{}
		err = json.Unmarshal(js, &ev)
		if err != nil {
			return err
		}
		err = s.handleCreateRoom(&ev, conn, token)
		if err != nil {
			return err
		}
	case EventTypeLeaveRoom:
		err = s.handleLeaveRoom(token)
		if err != nil {
			return err
		}
	case EventTypeStartGame:
		ev := EventStartGame{}
		err = json.Unmarshal(js, &ev)
		if err != nil {
			return err
		}
		err = s.handleStartGame(&ev, conn)
		if err != nil {
			return err
		}
	case EventTypeJoinRoom:
		ev := EventJoinRoom{}
		err = json.Unmarshal(js, &ev)
		if err != nil {
			return err
		}
		err = s.handleJoinRoom(&ev, conn, token)
		if err != nil {
			return err
		}

	}
	return conn.WriteJSON(event)
}


func (s *Server) handleCreateRoom(event *EventCreateRoom, conn *websocket.Conn, token jwt.Claims) error {
	player := werewolves.NewPlayer(token.Username, token.Username)
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


func (s *Server) handleLeaveRoom(token jwt.Claims) error {
	player, ok := s.Players[token.Username]
	if !ok {
		return s.serverError(PlayerNotFoundErr, EventTypeLeaveRoom)
	}
	if player.Room() != nil {
		delete(s.PlayersRoom, token.Username)
		return player.Room().RemovePlayer(player.ID)
	}
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

func (s *Server) handleJoinRoom(event *EventJoinRoom, conn *websocket.Conn, token jwt.Claims) error {
	_, ok := s.PlayersRoom[token.Username]
	if ok {
		return s.serverError(AlreadyInRoomErr ,EventTypeJoinRoom)
	}
	room, ok := s.Rooms[event.RoomID]
	if !ok{
		return s.serverError(RoomNotExistsErr ,EventTypeJoinRoom)
	}
	if room.Started(){
		return s.serverError(RoomStartedErr ,EventTypeJoinRoom)
	}
	player := werewolves.NewPlayer(token.Username, token.Username)
	room.AddPlayer(player)
	s.PlayersRoom[token.Username] = room.ID
	p, ok := s.Players[token.Username]
	if !ok {
		return s.serverError(PlayerNotFoundErr, EventTypeJoinRoom)
	}
	p.SetRoom(room)
	p.SetChar(player)
	go func() {
		for {
			select {
			case value, ok := <-player.Update:
				if !ok{
					return
				}
				err := conn.WriteJSON(value)
				if err != nil {
					return
				}
			}
		}
	}()
	conn.WriteJSON(player)
	return nil
}