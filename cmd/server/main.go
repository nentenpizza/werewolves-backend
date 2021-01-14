package main

import (
	"log"
	"net/http"

	"github.com/nentenpizza/werewolves/server"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	server := server.New()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e := echo.New()
	e.GET("/ws", func() echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			wsEndpoint(c, server)
			return nil
		})
	}())

	e.Logger.Fatal(e.Start(":7070"))
}

func wsEndpoint(c echo.Context, server *server.Server) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	go server.WsReader(conn)
	return nil
}
