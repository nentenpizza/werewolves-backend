package websocket

import "github.com/nentenpizza/werewolves/wserver"

func (s *handler) OnDisconnect(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() == nil{
		s.c.Delete(client.Token.Username)
		return nil
	}
	client.Room().BroadcastEvent(Event{EventTypeDisconnected, EventPlayerID{PlayerID: client.Token.Username}})
	return nil
}