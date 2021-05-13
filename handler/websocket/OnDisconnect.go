package websocket

import "github.com/nentenpizza/werewolves/wserver"

func (h *handler) OnDisconnect(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() == nil {
		h.c.Delete(client.Token.Username)
		return nil
	}
	if room := client.Room(); room != nil {
		if !room.Started() {
			err := room.RemovePlayer(client.ID)
			if err != nil {
				return err
			}
			h.c.Delete(client.Token.Username)
			room.BroadcastEvent(Event{Type: EventTypeDisconnected,
				Data: EventPlayerID{PlayerID: client.Token.Username}})
			if len(room.Players) < 1 {
				h.deleteRoom(room.ID)
			}
		}
	}

	return nil
}
