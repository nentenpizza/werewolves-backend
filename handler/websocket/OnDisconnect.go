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
		err := room.RemovePlayer(client.ID)
		if err != nil {
			return err
		}
		room.BroadcastEvent(Event{Type: EventTypeDisconnected,
			Data: EventPlayerID{PlayerID: client.Token.Username}})
	}

	return nil
}
