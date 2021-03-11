package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
	"math/rand"
)

func (h *handler) endGame(e *werewolves.RoomResult, room *werewolves.Room) error {
	var loseGroup map[string]*werewolves.Player
	wonGroup, ok := room.Groups[e.WonGroup]
	if e.WonGroup == "wolves" {
		loseGroup, ok = room.Groups["peaceful"]
	} else {
		loseGroup, ok = room.Groups["wolves"]
	}
	if ok {
		for _, p := range wonGroup {
			user, err := h.db.Users.ByUsername(p.Name)
			if err != nil {
				continue
			}
			xp := rand.Intn(1000)
			user.XP += int64(xp)
			user.Wins++
			p.Update <- &Event{Type: EventTypeEndGame, Data: EventEndGame{WonGroup: wonGroup, LoseGroup: loseGroup, XP: xp}}
			err = h.db.Users.Update(user)
			if err != nil {
				return err
			}
			client := h.c.ByID(p.ID)
			client.LeaveRoom()
		}

		for _, p := range loseGroup {
			user, err := h.db.Users.ByUsername(p.Name)
			if err != nil {
				continue
			}
			user.Losses++
			p.Update <- &Event{Type: EventTypeEndGame, Data: EventEndGame{WonGroup: wonGroup, LoseGroup: loseGroup, XP: 0}}
			err = h.db.Users.Update(user)
			if err != nil {
				return err
			}
			client := h.c.ByID(p.ID)
			client.LeaveRoom()
		}
	}
	h.deleteRoom(room.ID)
	return nil
}
