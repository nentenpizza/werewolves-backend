package websocket

import (
	"github.com/google/uuid"
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) OnJoinRoom(ctx *wserver.Context) error {
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
	room := h.r.Read(event.RoomID)
	if room == nil {
		return RoomNotExistsErr
	}
	if room.Started() {
		return RoomStartedErr
	}

	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)
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

func (h *handler) OnStartGame(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	log.Println(client.Room())
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
	h.broadcastToClients(&Event{Type: EventTypeRoomDeleted, Data: &EventRoomDeleted{RoomID: room.ID}})
	go func() {
		select {
		case e := <-room.NotifyDone:
			h.endGame(e, room)
			return
		}
	}()
	return nil
}

func (h *handler) OnCreateRoom(ctx *wserver.Context) error {
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

	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)

	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, client.Token.Username)
	h.r.Write(room.ID, room)

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
		h.c.Lock()
		for _, c := range h.c.clients {
			c.WriteJSON(&Event{Type: EventTypeRoomCreated, Data: EventNewRoomCreated{Room: room}})
		}
		h.c.Unlock()
	}()
	return nil
}

func (h *handler) OnLeaveRoom(ctx *wserver.Context) error {
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
			h.deleteRoom(room.ID)
			h.broadcastToClients(&Event{Type: EventTypeRoomDeleted, Data: EventRoomDeleted{RoomID: room.ID}})
		}
	}
	client.WriteJSON(ctx.Update)
	return nil
}

func (h *handler) deleteRoom(roomID string) {
	h.r.Delete(roomID)
}

func (h handler) broadcastToClients(e interface{}) {
	go func() {
		h.c.Lock()
		for _, c := range h.c.clients {
			c.WriteJSON(e)
		}
		h.c.Unlock()
	}()
}
