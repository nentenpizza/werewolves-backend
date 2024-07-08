package werewolves

import (
	"github.com/nentenpizza/werewolves/wserver"
)

func (g *game) OnVote(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() == nil {
		return NotInRoomRoom
	}
	if client.Player == nil {
		return PlayerNotFoundErr
	}
	event := &EventPlayerID{}
	if err := ctx.Bind(event); err != nil {
		return err
	}
	return g.Vote(client.Player, event.PlayerID)
}
