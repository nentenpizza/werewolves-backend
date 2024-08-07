package werewolves

import "github.com/nentenpizza/werewolves/pkg/wserver"

func (g *game) OnDisconnect(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	if client.Room() == nil {
		g.c.Delete(client.Token.Username)
		return nil
	}
	if room := client.Room(); room != nil {
		if !room.Started() {
			err := room.RemovePlayer(client.ID)
			if err != nil {
				return err
			}
			g.c.Delete(client.Token.Username)
			room.BroadcastEvent(Event{Type: EventTypeDisconnected,
				Data: EventPlayerID{PlayerID: client.Token.Username}})
			if len(room.Players) < 1 {
				g.deleteRoom(room.ID)
			}
		}
	}

	friends, err := g.friends.UserFriends(client.Token.ID)
	if err != nil {
		return err
	}

	online := make([]string, 0)

	for _, f := range friends {
		c := g.c.Read(f.Username)
		if c != nil {
			online = append(online, f.Username)
			c.conn.WriteJSON(Event{Type: EventTypeFriendLoggedOut,
				Data: EventUsername{client.Token.Username},
			})
		}
	}

	return nil
}
