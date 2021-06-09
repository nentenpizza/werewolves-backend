package websocket

import "github.com/nentenpizza/werewolves/wserver"

func (h *handler) OnDisconnect(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	go func() {
		friends, err := h.friends.UserFriends(client.Token.ID)
		if err != nil {
			return
		}

		online := make([]string, 0)

		for _, f := range friends {
			c := h.c.Read(f.Username)
			if c != nil {
				online = append(online, f.Username)
				c.conn.WriteJSON(Event{Type: EventTypeFriendLoggedOut,
					Data: EventUsername{client.Token.Username},
				})
			}
		}

	}()

	if client.Room() == nil {
		h.c.Delete(client.Token.Username)
		return nil
	}
	if room := client.Room(); room != nil {
		if !room.Started() {
			err := room.RemovePlayer(client.ID)
			if err != nil {
				return err
			}
			h.c.Delete(client.Token.Username)
			room.BroadcastEvent(Event{Type: EventTypeDisconnected,
				Data: EventPlayerID{PlayerID: client.Token.Username}})
			if len(room.Players) < 1 {
				h.deleteRoom(room.ID)
			}
		}
	}

	return nil
}
