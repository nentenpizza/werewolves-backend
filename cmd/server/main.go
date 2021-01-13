package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e := echo.New()
	e.GET("/ws", wsEndpoint)
}

func wsEndpoint(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	for {

	}
	return nil
}
