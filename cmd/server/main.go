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
			err := s.WsEndpoint(c)
			return err
		})
	}())

	e.Logger.Fatal(e.Start(":7070"))
}
