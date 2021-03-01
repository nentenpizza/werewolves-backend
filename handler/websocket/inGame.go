package websocket

import (
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) OnVote(ctx wserver.Context) error  {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() == nil{
		return NotInRoomRoom
	}

	event := &EventPlayerID{}
	err := h.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}
	log.Println("vote", ctx.EventType(), event.PlayerID)
	action := client.Player.Vote(event.PlayerID)
	err = client.Room().Perform(action)
	if err != nil {
		return NotAllowedErr
	}
	return nil
}