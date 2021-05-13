package websocket

import (
	"github.com/nentenpizza/werewolves/wserver"
)

func (h *handler) OnConnect(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client != nil {
		if len(client.Unreached) > 1 {
			client.SendUnreached()
		}
		if client.Room() == nil {
			client.WriteJSON(&Event{Type: EventTypeNotInGame})
		} else {
			Logger.WithField("client_room", client.room).Info("reconnected to room")
		}
	}
	return ctx.Conn.WriteJSON(Event{Type: EventTypeAllRooms, Data: EventAllRooms{h.r}})
}
