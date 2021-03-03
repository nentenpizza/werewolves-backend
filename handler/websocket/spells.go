package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
)

func (h *handler) OnSkill(ctx wserver.Context) error  {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil {
		switch char := client.Player.Character.(type) {
		case *werewolves.Constable:
			e := &TargetedEvent{}
			if err := ctx.Bind(e); err != nil {
				return err
			}
			p, ok := client.Room().Players[e.TargetID]
			if !ok {
				return PlayerNotFoundErr
			}
			action := char.Shoot(p)
			err := client.Room().Perform(action)
			if err != nil {
				return NotAllowedErr
			}
		}
	}
	return nil
}