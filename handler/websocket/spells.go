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
	if room := client.Room(); room != nil {

		var action werewolves.Action

		switch char := client.Player.Character.(type) {
		case *werewolves.Constable:
			e := TargetedEvent{}
			if err := ctx.Bind(&e); err != nil {
				return err
			}
			p, ok := room.Players[e.TargetID]
			if !ok {
				return PlayerNotFoundErr
			}
			action = char.Shoot(p)

		case *werewolves.Doctor:
			e := TargetedEvent{}
			if err := ctx.Bind(&e); err != nil {
				return err
			}
			p, ok := room.Players[e.TargetID]
			if !ok {
				return PlayerNotFoundErr
			}
			action = char.Heal(p)

		case *werewolves.AlphaWerewolf:
			if room.State != werewolves.Night{
				return NotAllowedErr
			}
			e := TargetedEvent{}
			if err := ctx.Bind(&e); err != nil {
				return err
			}
			p, ok := room.Players[e.TargetID]
			if !ok {
				return PlayerNotFoundErr
			}
			err := client.WriteJSON(Event{Type: EventTypeRevealRole, Data: EventRevealRole{p.Role,p.ID}})
			return err
		}

		err := client.Room().Perform(action)
		if err != nil {
			return NotAllowedErr
		}

	}
	return nil
}