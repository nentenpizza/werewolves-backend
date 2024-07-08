package werewolves

import (
	"encoding/json"
	"sync"

	"github.com/nentenpizza/werewolves/pkg/werewolves"
	"github.com/nentenpizza/werewolves/pkg/wserver"
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

func NewRooms() *Rooms {
	return &Rooms{rooms: make(map[string]*werewolves.Room)}
}

func (g *game) OnListRooms(c *wserver.Context) error {
	c.Update.Data = EventAllRooms{Rooms: g.r}
	return c.Conn.WriteJSON(c.Update)
}
