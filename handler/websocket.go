package handler

import (
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/werewolves"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) WsEndpoint(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	go s.WsReader(conn)
	return nil
}

func (s *Server) AllRooms(c echo.Context) error {
	return c.JSON(http.StatusOK, s.Rooms)
}

func (s *Server) JoinRoom(c echo.Context) error {
	var form struct {
		RoomID string `json:"room_id"`
	}
	if err := c.Bind(&form); err != nil{
		return err
	}
	token := jwt.From(c.Get("user"))

	_, ok := s.PlayersRoom[token.Username]
	if ok {
		return c.JSON(http.StatusConflict, app.Err("you already in room"))
	}
	room, ok := s.Rooms[form.RoomID]
	if !ok{
		return c.JSON(http.StatusForbidden, app.Err("room does not exists"))
	}
	if room.Started(){
		return c.JSON(http.StatusConflict, app.Err("room already started"))
	}
	player := werewolves.NewPlayer(token.Username, token.Username)
	room.AddPlayer(player)
	s.PlayersRoom[token.Username] = room.ID
	return c.JSON(200, app.Ok())
}