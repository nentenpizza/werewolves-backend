package transport

import (
	"github.com/google/uuid"
	werewolves2 "github.com/nentenpizza/werewolves/game/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
)

func (g *game) OnJoinRoom(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	event := &EventRoomPlayer{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}
	//c := s.c.Read(client.Token.Username)

	if client.Room() != nil {
		return AlreadyInRoomErr
	}
	room := g.r.Read(event.RoomID)
	if room == nil {
		return RoomNotExistsErr
	}
	if room.Started() {
		return RoomStartedErr
	}

	player := werewolves2.NewPlayer(client.Token.Username, client.Token.Username)
	room.AddPlayer(player)

	event.PlayerID = client.Token.Username
	ctx.Update.Data = event

	client.SetRoom(room)
	client.SetChar(player)

	go client.ListenRoom()

	room.BroadcastEvent(ctx.Update)

	client.WriteJSON(player)
	return nil
}

func (g *game) OnStartGame(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	room := client.Room()
	if room == nil {
		return NotInRoomRoom
	}
	if client.Token.Username != room.Owner {
		return NotAllowedErr
	}
	err := room.Start()
	if err != nil {
		return RoomStartErr
	}
	g.broadcastToClients(&Event{Type: EventTypeRoomDeleted, Data: &EventRoomDeleted{RoomID: room.ID}})
	go func() {
		select {
		case e := <-room.NotifyDone:
			g.endGame(e, room)
			return
		}
	}()
	return nil
}

func (g *game) OnCreateRoom(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil {
		return AlreadyInRoomErr
	}

	event := &EventCreateRoom{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}

	player := werewolves2.NewPlayer(client.Token.Username, client.Token.Username)

	roomID := uuid.New().String()
	room := werewolves2.NewRoom(roomID, event.RoomName, werewolves2.Players{}, event.Settings, client.Token.Username)
	g.r.Write(room.ID, room)

	err = room.AddPlayer(player)
	if err != nil {
		return err
	}

	client.SetRoom(room)
	client.SetChar(player)

	go client.ListenRoom()

	err = client.WriteJSON(player)
	if err != nil {
		return err
	}

	go func() {
		g.c.Lock()
		for _, c := range g.c.clients {
			c.WriteJSON(&Event{Type: EventTypeRoomCreated, Data: EventNewRoomCreated{Room: room}})
		}
		g.c.Unlock()
	}()
	return nil
}

func (g *game) OnLeaveRoom(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if room := client.Room(); room != nil {
		err := room.RemovePlayer(client.ID)
		if err != nil {
			return err
		}
		ctx.Update.Data = EventRoomPlayer{PlayerID: client.Token.Username}
		room.BroadcastEvent(ctx.Update)
		client.SetRoom(nil)
		client.quit <- true
		if len(room.Players) <= 0 {
			g.deleteRoom(room.ID)
			g.broadcastToClients(&Event{Type: EventTypeRoomDeleted, Data: EventRoomDeleted{RoomID: room.ID}})
		}
	}
	client.WriteJSON(ctx.Update)
	g.c.Delete(client.Token.Username)
	return nil
}

func (g *game) deleteRoom(roomID string) {
	g.r.Delete(roomID)
}

func (g game) broadcastToClients(e interface{}) {
	go func() {
		g.c.Lock()
		for _, c := range g.c.clients {
			c.WriteJSON(e)
		}
		g.c.Unlock()
	}()
}
