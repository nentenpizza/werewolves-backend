package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/handler"
)

func main() {
	s := handler.NewServer()

	e := echo.New()
	e.GET("/ws", func() echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			s.WsEndpoint(c)
			return nil
		})
	}())

	e.Logger.Fatal(e.Start(":7070"))
}
