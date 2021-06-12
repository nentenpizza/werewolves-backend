package transport

import (
	"github.com/nentenpizza/werewolves/game/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"time"
)

var FloodWaitDuration = 2 * time.Second
var EmojiWaitDuration = 2 * time.Second

func (g *game) OnMessage(ctx *wserver.Context) error {
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
			if time.Now().Sub(client.FloodWait) < FloodWaitDuration {
				return client.WriteJSON(Event{
					EventTypeFloodWait,
					EventFloodWait{
						int64((FloodWaitDuration-time.Since(client.FloodWait))/time.Second) + 1},
				})
			}
			if len([]rune(event.Text)) > 160 {
				return NotAllowedErr
			}
			if client.Player.Character != nil {
				if client.Player.Character.HP() > 0 {
					if client.Room().State != werewolves.Night {
						client.Room().BroadcastEvent(ctx.Update)
					} else {
						if client.Room().InGroup("wolves", client.Player.ID) {
							client.Room().BroadcastTo("wolves", ctx.Update)
						}
					}
				} else {
					client.Room().BroadcastTo("dead", ctx.Update)
				}

				client.FloodWait = time.Now()
			} else {
				client.Room().BroadcastEvent(ctx.Update)
			}
		}
	}
	return nil
}

func (g *game) OnEmote(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client != nil {
		if time.Now().Sub(client.EmojiWait) < EmojiWaitDuration {
			return client.WriteJSON(Event{
				EventTypeFloodWait,
				EventFloodWait{
					int64((EmojiWaitDuration-time.Since(client.EmojiWait))/time.Second) + 1},
			})
		}

		event := &EmoteEvent{}
		err := ctx.Bind(event)
		if err != nil {
			return err
		}
		if room := client.Room(); room != nil {
			room.BroadcastEvent(Event{EventTypeSendEmote, EmoteEvent{Emote: event.Emote, FromID: client.Token.Username}})

			client.EmojiWait = time.Now()
		}
	}

	return nil
}
