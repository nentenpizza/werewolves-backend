package werewolves

import (
	"errors"
	"math"
	"math/rand"

	werewolves "github.com/nentenpizza/werewolves/pkg/werewolves"
)

func (g *game) endGame(e *werewolves.RoomResult, room *werewolves.Room) error {
	var loseGroup map[string]*werewolves.Player
	wonGroup, ok := room.Groups[e.WonGroup]
	if !ok {
		return errors.New("websocket: wonGroup not found")
	}
	if e.WonGroup == "wolves" {
		loseGroup, ok = room.Groups["peaceful"]
	} else {
		loseGroup, ok = room.Groups["wolves"]
	}
	if !ok {
		return errors.New("websocket: wonGroup not found")
	}
	for _, p := range wonGroup {
		user, err := g.db.Users.ByUsername(p.Name)
		if err != nil {
			continue
		}
		xp := int(math.Max(500, float64(rand.Intn(1000))))
		user.XP += int64(xp)
		user.Wins++

		err = g.db.Users.Update(user)
		if err != nil {
			return err
		}

		client := g.c.ByID(p.ID)
		if client != nil {
			p.Update <- &Event{Type: EventTypeEndGame, Data: EventEndGame{WonGroup: wonGroup, LoseGroup: loseGroup, XP: xp}}
			g.c.Delete(client.ID)
		}

	}

	for _, p := range loseGroup {
		user, err := g.db.Users.ByUsername(p.Name)
		if err != nil {
			continue
		}
		user.Losses++
		err = g.db.Users.Update(user)
		if err != nil {
			return err
		}

		client := g.c.ByID(p.ID)
		if client != nil {
			p.Update <- &Event{Type: EventTypeEndGame, Data: EventEndGame{WonGroup: wonGroup, LoseGroup: loseGroup, XP: 0}}
			g.c.Delete(client.ID)
		}
	}
	g.r.Delete(room.ID)
	g.broadcastToClients(&Event{Type: EventTypeRoomDeleted, Data: EventRoomDeleted{RoomID: room.ID}})
	return nil
}
