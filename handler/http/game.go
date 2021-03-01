package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/werewolves"
	"time"
)

type GameService struct {
	handler
}

func (s GameService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/phase_length", s.GetPhaseLength)
}

func (s GameService) GetPhaseLength(c echo.Context) error {
	return c.JSON(200, echo.Map{"length": werewolves.PhaseLength / time.Second})
}