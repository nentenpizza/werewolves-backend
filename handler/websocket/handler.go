package websocket

import (
	"encoding/json"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/werewolves"
	"math/rand"
	"strconv"
	"sync"
)

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
	DB      *storage.DB
	Rooms   *Rooms
	Clients *Clients
	Secret  []byte
}

type handler struct {
	db *storage.DB
	r  *Rooms
	c  *Clients
	s  []byte
}

func NewHandler(h Handler) *handler {
	for i := 0; i < 5; i++ {
		room := werewolves.NewRoom(strconv.Itoa(rand.Intn(100)), strconv.Itoa(rand.Intn(100)), werewolves.Players{}, werewolves.Settings{}, strconv.Itoa(rand.Intn(100)))
		h.Rooms.Write(room.ID, room)
	}
	return &handler{
		db: h.DB,
		r:  h.Rooms,
		c:  h.Clients,
		s:  h.Secret,
	}
}
