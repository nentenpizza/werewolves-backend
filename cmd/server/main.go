package main

import (
	"github.com/nentenpizza/werewolves/server"

	"github.com/labstack/echo/v4"
)

func main() {
	s := server.New()

	e := echo.New()
	e.GET("/ws", func() echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			s.WsEndpoint(c)
			return nil
		})
	}())

	e.Logger.Fatal(e.Start(":7070"))
}
