package werewolves

import (
	"encoding/json"
	"sync"

	"github.com/nentenpizza/werewolves/internal/service"
	"github.com/nentenpizza/werewolves/internal/storage"
	"github.com/nentenpizza/werewolves/pkg/wserver"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.TextFormatter{})
}

type Clients struct {
	clients map[string]*Client
	sync.Mutex
}

func (m *Clients) Write(key string, value *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[key] = value
}

func (m *Clients) Read(key string) *Client {
	m.Lock()
	defer m.Unlock()
	return m.clients[key]
}

func (m *Clients) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.clients, key)
}

func (m *Clients) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(m.clients)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (m *Clients) ByID(key string) *Client {
	m.Lock()
	defer m.Unlock()
	return m.clients[key]
}

func NewClients(c map[string]*Client) *Clients {
	return &Clients{clients: c}
}

type Game struct {
	DB             *storage.DB
	Rooms          *Rooms
	Clients        *Clients
	Secret         []byte
	FriendsService service.FriendService
}

type game struct {
	db      *storage.DB
	r       *Rooms
	c       *Clients
	s       []byte
	friends service.FriendService
}

func NewGame(h Game) *game {
	return &game{
		db:      h.DB,
		r:       h.Rooms,
		c:       h.Clients,
		s:       h.Secret,
		friends: h.FriendsService,
	}
}

func (g game) Register(server *wserver.Server) {
	server.Use(g.WebsocketJWT, g.Logger)

	server.Handle(EventTypeCreateRoom, g.OnCreateRoom)
	server.Handle(EventTypeJoinRoom, g.OnJoinRoom)
	server.Handle(EventTypeLeaveRoom, g.OnLeaveRoom)
	server.Handle(EventTypeStartGame, g.OnStartGame)

	server.Handle(wserver.OnConnect, g.OnConnect)
	server.Handle(wserver.OnDisconnect, g.OnDisconnect)

	server.Handle(EventTypeSendMessage, g.OnMessage)
	server.Handle(EventTypeVote, g.OnVote)
	server.Handle(EventTypeUseSkill, g.OnSkill)
	server.Handle(EventTypeSendEmote, g.OnEmote)
	server.Handle(EventTypeAllRooms, g.OnListRooms)
}
