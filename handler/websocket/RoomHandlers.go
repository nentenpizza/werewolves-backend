package websocket

import (
	"github.com/google/uuid"
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (s *handler) OnJoinRoom(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	event := &EventRoomPlayer{}
	err := s.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}
	//c := s.c.Read(client.Token.Username)

	if client.Room() != nil {
		return AlreadyInRoomErr
	}
	room := s.r.Read(event.RoomID)
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

	s.c.Write(client.Token.Username, client)

	go client.ListenRoom()

	room.BroadcastEvent(ctx.Update)

	client.WriteJSON(player)
	return nil
}


func (s *handler) OnStartGame(ctx wserver.Context) error {
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

func (s *handler) OnCreateRoom(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil{
		return AlreadyInRoomErr
	}

	event := &EventCreateRoom{}
	err := s.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}

	player := werewolves.NewPlayer(client.Token.Username, client.Token.Username)

	roomID := uuid.New().String()
	room := werewolves.NewRoom(roomID, event.RoomName, werewolves.Players{}, event.Settings, client.Token.Username)
	s.r.Write(room.ID,room)

	err = room.AddPlayer(player)
	if err != nil {
		return err
	}

	client.SetRoom(room)
	client.SetChar(player)

	s.c.Write(client.Token.Username, client)

	go client.ListenRoom()

	err = client.WriteJSON(player)
	if err != nil {
		return err
	}

	go func(){
		s.c.Lock()
		for _, c := range s.c.clients{
			c.WriteJSON(Event{EventTypeRoomCreated, EventNewRoomCreated{Room: room}})
		}
		s.c.Unlock()
	}()
	return nil
}

func (s *handler) OnLeaveRoom(ctx wserver.Context) error {
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
		s.c.Delete(client.Token.Username)
	}
	client.WriteJSON(ctx.Update)
	return nil
}

func (s *handler) OnMessage(ctx wserver.Context) error {
	client := ctx.Get("client").(*Client)

	event := &MessageEvent{}
	err := s.mapToEvent(event, ctx.Data())
	if err != nil {
		return err
	}

	if client != nil {
		if client.Room() != nil {
			event.Username = client.Token.Username
			ctx.Update.Data = event
			client.Room().BroadcastEvent(ctx.Update)
		}
	}
	return nil
}