package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) OnVote(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() == nil {
		return NotInRoomRoom
	}
	if client.Player == nil{
		return PlayerNotFoundErr
	}
	var vCount uint8 = 1
	if client.Player.Role == "AlphaWerewolf" && client.Room().State == werewolves.Night{
		vCount = 2
	}
	event := &EventPlayerID{}
	if err := ctx.Bind(event); err != nil {
		return err
	}
	log.Println(event.PlayerID)
	action := client.Player.Vote(event.PlayerID, vCount)
	if client.Room().State == werewolves.Night {
		if err := client.Room().Perform(action, "wolves"); err != nil {
			return NotAllowedErr
		}
		log.Println("night vote", ctx.EventType(), event.PlayerID)
	} else {
		if err := client.Room().Perform(action); err != nil {
			return NotAllowedErr
		}
		log.Println("day vote", ctx.EventType(), event.PlayerID)
	}
	return nil
}
