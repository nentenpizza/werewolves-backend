package websocket

import (
	"github.com/google/uuid"
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) OnJoinRoom(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	event := &EventRoomPlayer{}
	err := h.mapToEvent(event, ctx.Data())
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
	if room.Started(){
		return RoomStartedErr
	}

	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)
	room.AddPlayer(player)

	event.PlayerID = client.Token.Username
	ctx.Update.Data = event

	client.SetRoom(room)
	client.SetChar(player)

	h.c.Write(client.Token.Username, client)

	go client.ListenRoom()

	room.BroadcastEvent(ctx.Update)

	client.WriteJSON(player)
	return nil
}


func (h *handler) OnStartGame(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	log.Println(client.Room())
	room := client.Room()
	if room == nil{
		return NotInRoomRoom
	}
	if client.Token.Username != room.Owner {
		return NotAllowedErr
	}
	err := room.Start()
	if err != nil {
		return RoomStartErr
	}
	return nil
}

func (h *handler) OnCreateRoom(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil{
		return AlreadyInRoomErr
	}

	event := &EventCreateRoom{}
	err := h.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}

	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)

	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, client.Token.Username)
	h.r.Write(room.ID,room)

	err = room.AddPlayer(player)
	if err != nil {
		return err
	}

	client.SetRoom(room)
	client.SetChar(player)

	h.c.Write(client.Token.Username, client)

	go client.ListenRoom()

	err = client.WriteJSON(player)
	if err != nil {
		return err
	}

	go func(){
		h.c.Lock()
		for _, c := range h.c.clients{
			c.WriteJSON(Event{Type: EventTypeRoomCreated, Data: EventNewRoomCreated{Room: room}})
		}
		h.c.Unlock()
	}()
	return nil
}

func (h *handler) OnLeaveRoom(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil {
		err := client.Room().RemovePlayer(client.ID)
		if err != nil {
			return err
		}
		ctx.Update.Data = EventRoomPlayer{PlayerID: client.Token.Username}
		client.Room().BroadcastEvent(ctx.Update)
		client.SetRoom(nil)
		client.quit <- true
		h.c.Delete(client.Token.Username)
	}
	client.WriteJSON(ctx.Update)
	return nil
}

func (h *handler) OnMessage(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)

	event := &MessageEvent{}
	err := h.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}

	if client != nil {
		if client.Room() != nil {
			event.Username = client.Token.Username
			ctx.Update.Data = event
			if client.Player.Character.HP() <= 0 {
				client.Room().BroadcastTo("dead",ctx.Update)
				return nil
			}
			if client.Room().State != werewolves.Night {
				client.Room().BroadcastEvent(ctx.Update)
			}else{
				if client.Room().InGroup("wolves", client.Player.ID){
					client.Room().BroadcastTo("wolves",ctx.Update)
				}
			}
		}
	}
	return nil
}	