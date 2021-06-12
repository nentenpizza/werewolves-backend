package transport

import (
	"encoding/json"
	"github.com/nentenpizza/werewolves/game/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"sync"
)

type Rooms struct {
	rooms map[string]*werewolves.Room
	sync.Mutex
}

func (m *Rooms) Write(key string, value *werewolves.Room) {
	m.Lock()
	defer m.Unlock()
	m.rooms[key] = value
}

func (m *Rooms) Read(key string) *werewolves.Room {
	m.Lock()
	defer m.Unlock()
	return m.rooms[key]
}

func (m *Rooms) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.rooms, key)
}

func (m *Rooms) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(m.rooms)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func NewRooms(m map[string]*werewolves.Room) *Rooms {
	return &Rooms{rooms: m}
}

func (g *game) OnListRooms(c *wserver.Context) error {
	c.Update.Data = EventAllRooms{Rooms: g.r}
	return c.Conn.WriteJSON(c.Update)
}
