package transport

import (
	"github.com/nentenpizza/werewolves/wserver"
	"log"
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
		go func() {
			friends, err := g.friends.UserFriends(client.Token.ID)
			if err != nil {
				return
			}

			online := make([]string, 0)

			for _, f := range friends {
				c := g.c.Read(f.Username)
				if c != nil {
					online = append(online, f.Username)
					c.conn.WriteJSON(Event{Type: EventTypeFriendLoggedIn,
						Data: EventUsername{client.Token.Username},
					})
				}
			}

			err = ctx.Conn.WriteJSON(Event{Type: EventTypeFriendsOnlineInfo,
				Data: EventFriendsOnlineInfo{online},
			})
			if err != nil {
				log.Println(err)
			}
		}()
	}

	return ctx.Conn.WriteJSON(Event{Type: EventTypeAllRooms, Data: EventAllRooms{g.r}})
}
