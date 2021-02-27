package websocket

import (
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (s *handler) OnConnect(ctx wserver.Context) error  {
	client := ctx.Get("client").(*Client)
	if client != nil{
		if len(client.Unreached) > 1{
			client.SendUnreached()
		}
	}
	log.Println("Connect")
	return ctx.Conn.WriteJSON(Event{EventTypeAllRooms, EventAllRooms{s.r}})
}