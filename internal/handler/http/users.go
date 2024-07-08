package http

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/pkg/jwt"
)

type UsersEndpointGroup struct {
	handler
}

func (s UsersEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/me", s.Me)
	g.POST("/user", s.GetUser)
}

func (s UsersEndpointGroup) Me(c echo.Context) error {
	token := jwt.From(c.Get("user"))

	user, err := s.db.Users.ByUsername(token.Username)
	if err != nil {

		return err
	}
	return c.JSON(200, user)
}

func (s UsersEndpointGroup) GetUser(c echo.Context) error {
	var form struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}

	user, err := s.db.Users.ByUsername(form.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return c.JSON(http.StatusNotFound, app.Err("user not found"))
	}
	return c.JSON(200, user)
}
