package websocket

import (
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
)

func (h *handler) OnMessage(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)

	event := &MessageEvent{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}

	if client != nil {
		if client.Room() != nil {
			event.Username = client.Token.Username
			ctx.Update.Data = event
			if client.Player.Character.HP() <= 0 {
				client.Room().BroadcastTo("dead", ctx.Update)
				return nil
			}
			if client.Room().State != werewolves.Night {
				client.Room().BroadcastEvent(ctx.Update)
			} else {
				if client.Room().InGroup("wolves", client.Player.ID) {
					client.Room().BroadcastTo("wolves", ctx.Update)
				}
			}
		}
	}
	return nil
}

func (h *handler) OnEmote(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)

	event := &EmoteEvent{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}
	if room := client.Room(); room != nil {
		room.BroadcastEvent(Event{EventTypeSendEmote, EmoteEvent{Emote: event.Emote, FromID: client.Token.Username}})
	}

	return nil
}
