package handler

import (
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
