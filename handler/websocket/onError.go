package websocket

import (
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) OnError(err error, ctx wserver.Context) {
	//client, ok := ctx.Get("client").(*Client)
	//if ok {
	//	if err != nil{
	//		log.Println(client.Name, client.Role, err)
	//
	//
	//	e, k := err.(*ServerError)
	//	if k {
	//		client.WriteJSON(Event{Type: ctx.EventType(), Data: ctx.Data(), Error: *e})
	//	}
	//	}
	//	return
	//}
	if err != nil{
		log.Println(err, ctx)
	}
}