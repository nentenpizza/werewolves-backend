package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
	"math/rand"
)

func (h *handler) endGame(e *werewolves.RoomResult, room *werewolves.Room) {
	var loseGroup map[string]*werewolves.Player
	wonGroup, ok := room.Groups[e.WonGroup]
	if ok {
		for _, p := range wonGroup {
			user, err := h.db.Users.ByUsername(p.Name)
			if err != nil {
				continue
			}
			user.XP += int64(rand.Intn(1000))
			user.Wins++
			err = h.db.Users.Update(user)
		}
	}
	if e.WonGroup == "wolves" {
		loseGroup, ok = room.Groups["peaceful"]
	} else {
		loseGroup, ok = room.Groups["wolves"]
	}
	if ok {
		for _, p := range loseGroup {
			user, err := h.db.Users.ByUsername(p.Name)
			if err != nil {
				continue
			}
			user.Losses++
			err = h.db.Users.Update(user)
		}
	}
	h.deleteRoom(room.ID)
	h.broadcastToClients(Event{Type: EventTypeEndGame, Data: EventEndGame{WonGroup: wonGroup, LoseGroup: loseGroup, XP: 228}})
	room = nil
	return
}
