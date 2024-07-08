package werewolves

import (
	"github.com/nentenpizza/werewolves/pkg/wserver"
)

func (g *game) OnConnect(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client != nil {
		if len(client.Unreached) > 1 {
			client.SendUnreached()
		}
		if client.Room() == nil {
			client.WriteJSON(&Event{Type: EventTypeNotInGame})
		} else {
			Logger.WithField("client_room", client.room).Info("reconnected to room")
		}

		friends, err := g.friends.UserFriends(client.Token.ID)
		if err != nil {
			return err
		}

		usersOnline := make([]string, 0)

		for _, f := range friends {
			c := g.c.Read(f.Username)
			if c != nil {
				usersOnline = append(usersOnline, f.Username)
				err = c.conn.WriteJSON(Event{Type: EventTypeFriendLoggedIn,
					Data: EventUsername{client.Token.Username},
				})
				if err != nil {
					Logger.WithField("client_id", client.ID).Error(err)
				}
			}
		}

		err = ctx.Conn.WriteJSON(Event{Type: EventTypeFriendsOnlineInfo,
			Data: EventFriendsOnlineInfo{usersOnline},
		})
		if err != nil {
			Logger.WithField("client_id", client.ID).Error(err)
		}
	}

	return ctx.Conn.WriteJSON(Event{Type: EventTypeAllRooms, Data: EventAllRooms{g.r}})
}
