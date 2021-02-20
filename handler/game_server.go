package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

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

const reconnectTime = 5 * time.Second

type Client struct {
	sync.Mutex
	conn *websocket.Conn
	*werewolves.Player
	room *werewolves.Room
	AFK bool
	Token jwt.Claims
	Unreached []interface{}
	quit chan bool
}

func (c *Client) Conn() *websocket.Conn {
	c.Lock()
	defer c.Unlock()
	return c.conn
}

func (c *Client) UpdateConn(conn *websocket.Conn){
	c.Lock()
	defer c.Unlock()
	c.conn = conn
}

func (c *Client) SetRoom(r *werewolves.Room){
	c.Lock()
	defer c.Unlock()
	c.room = r
}
func (c *Client) SetChar(plr *werewolves.Player){
	c.Lock()
	defer c.Unlock()
	c.Player = plr
}

func (c *Client) Room()*werewolves.Room{
	c.Lock()
	defer c.Unlock()
	return c.room
}
func (c *Client) Char() *werewolves.Player{
	c.Lock()
	defer c.Unlock()
	return c.Player
}

func (c *Client) ListenRoom(){
	for {
		if c.Player != nil {
			select {
			case value, ok := <-c.Player.Update:
				if ok {
					c.WriteJSON(value)
				} else {
					return
				}
			case <- c.quit:
				return
			}
		}else{
			return
		}
	}
}

func (c *Client) WriteJSON(i interface{}) error {
	c.Lock()
	defer c.Unlock()
	err := c.conn.WriteJSON(i)
	if err != nil{
		c.Unreached = append(c.Unreached, i)
		log.Println(c.Token.Username, "unreached", i)
	}
	return err
}

func (c *Client) SendUnreached(){
	if len(c.Unreached) > 0{
		for _, e := range c.Unreached{
			c.WriteJSON(e)
			log.Println("sent unreached to", c.Token.Username, "|", e)
		}
	}
	c.Unreached = make([]interface{}, 0)
}

// Server represents a game server which talks with game
type Server struct {
	handler
	Rooms       map[string]*werewolves.Room
	Clients map[string]*Client
	Secret []byte
	sync.Mutex
}

func (s *Server) REGISTER(h handler, g *echo.Group)  {
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
		Clients: make(map[string]*Client),
		Secret: secret,
	}
}

func (s *Server) WsReader(conn *websocket.Conn, token jwt.Claims) {
	client, ok := s.Clients[token.Username]
	if !ok {
		client = &Client{conn: conn, Token: token, AFK: false, quit: make(chan bool, 1)}
		s.Clients[token.Username] = client
	}else {
		log.Printf("reconnected %s", token.Username)
		client.AFK = false
		client.Token = token
		client.UpdateConn(conn)
		client.SendUnreached()
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Printf("waiting for %s", token.Username)
			client.AFK = true
			time.AfterFunc(reconnectTime, func(){
				if client.AFK == true {
					s.Disconnect(client)
				}
			})
			return
		}
		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(string(msg))
		err = s.HandleEvent(&ev, client)
		if err != nil {
			log.Println(err)
			client.WriteJSON(err)
		}
	}
}


func(s *Server) Disconnect(client *Client){
	client.quit<-true
	if client.Room() != nil{
		err := client.Room().RemovePlayer(client.Token.Username)
		log.Println(err)
	}
	delete(s.Clients ,client.Token.Username)
	log.Printf("disconnected %s", client.Token.Username)

}

func (s *Server) HandleEvent(event *Event, client *Client) error {
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
		err = s.handleCreateRoom(&ev, client)
		if err != nil {
			return err
		}
	case EventTypeLeaveRoom:
		err = s.handleLeaveRoom(event, client)
		if err != nil {
			return err
		}
	case EventTypeStartGame:
		err = s.handleStartGame(client)
		if err != nil {
			return err
		}
	case EventTypeJoinRoom:
		ev := EventJoinRoom{}
		err = json.Unmarshal(js, &ev)
		if err != nil {
			return err
		}
		err = s.handleJoinRoom(&ev,event, client)
		if err != nil {
			return err
		}

	}
	return client.WriteJSON(event)
}


func (s *Server) handleCreateRoom(event *EventCreateRoom,client *Client) error {
	c := s.clientByUsername(client.Token.Username)
	if c == nil{
		return s.serverError(PlayerNotFoundErr, EventTypeCreateRoom)
	}
	if c.Room() != nil{
		return s.serverError(AlreadyInRoomErr, EventTypeCreateRoom)
	}
	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)
	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, client.Token.Username)
	s.Rooms[room.ID] = room
	err := room.AddPlayer(player)
	if err != nil {
		return err
	}
	c.SetRoom(room)
	c.SetChar(player)
	go c.ListenRoom()
	err = client.WriteJSON(player)
	if err != nil {
		return err
	}



	return nil
}


func (s *Server) handleLeaveRoom(ev *Event, client *Client) error {
	if client.Room() != nil {
		err := client.Room().RemovePlayer(client.ID)
		ev.Data = EventLeaveRoom{PlayerID: client.Token.Username}
		client.Room().BroadcastEvent(ev)
		client.SetRoom(nil)
		client.quit <- true
		return err
	}
	return nil
}

func (s *Server) handleStartGame(client *Client) error {
	room := client.Room()
	if room == nil{
		return s.serverError(NotInRoomRoom, EventTypeStartGame)
	}
	if client.Token.Username != room.Owner {
		return s.serverError(NotAllowedErr, EventTypeStartGame)
	}
	err := room.Start()
	if err != nil {
		return s.serverError(RoomStartErr, EventTypeStartGame, err.Error())
	}
	return nil
}

func (s *Server) handleJoinRoom(event *EventJoinRoom, ev *Event, client *Client) error {
	c, ok := s.Clients[client.Token.Username]
	if !ok {
		return s.serverError(PlayerNotFoundErr ,EventTypeJoinRoom)
	}
	if c.Room() != nil {
		return s.serverError(AlreadyInRoomErr ,EventTypeJoinRoom)
	}
	room, ok := s.Rooms[event.RoomID]
	if !ok{
		return s.serverError(RoomNotExistsErr ,EventTypeJoinRoom)
	}
	if room.Started(){
		return s.serverError(RoomStartedErr ,EventTypeJoinRoom)
	}
	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)
	room.AddPlayer(player)
	if !ok {
		return s.serverError(PlayerNotFoundErr, EventTypeJoinRoom)
	}
	event.PlayerID = client.Token.Username
	ev.Data = event
	c.SetRoom(room)
	c.SetChar(player)
	go c.ListenRoom()
	room.BroadcastEvent(ev)

	client.WriteJSON(player)
	return nil
}

func (s *Server) clientByUsername(username string) *Client{
	c, _ := s.Clients[username]
	return c
}