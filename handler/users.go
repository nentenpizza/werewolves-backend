package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
)

type UsersService struct {
	handler
}

func (s UsersService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/me", s.Me)
}


func (s UsersService) Me(c echo.Context)error {
	token := jwt.From(c.Get("user"))

	user, err := s.db.Users.ByUsername(token.Username)
	if err != nil {

		return err
	}
	return c.JSON(200, user)
}
