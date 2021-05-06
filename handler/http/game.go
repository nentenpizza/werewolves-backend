package http

import (
	"github.com/labstack/echo/v4"
)

type GameService struct {
	handler
	PhaseLength int
}

func (s GameService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/phase_length", s.GetPhaseLength)
}

func (s GameService) GetPhaseLength(c echo.Context) error {
	return c.JSON(200, echo.Map{"length": s.PhaseLength})
}
