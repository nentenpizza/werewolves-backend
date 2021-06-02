package websocket

import (
	"encoding/json"
	"github.com/nentenpizza/werewolves/service"
	"github.com/nentenpizza/werewolves/storage"
	log "github.com/sirupsen/logrus"
	"sync"
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

type Handler struct {
	DB             *storage.DB
	Rooms          *Rooms
	Clients        *Clients
	Secret         []byte
	FriendsService service.FriendService
}

type handler struct {
	db      *storage.DB
	r       *Rooms
	c       *Clients
	s       []byte
	friends service.FriendService
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:      h.DB,
		r:       h.Rooms,
		c:       h.Clients,
		s:       h.Secret,
		friends: h.FriendsService,
	}
}
