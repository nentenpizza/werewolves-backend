package server

import (
	j "github.com/dgrijalva/jwt-go"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) WsEndpoint(c echo.Context) error {
	t := c.Param("token")
	if t == ""{
		return c.JSON(http.StatusBadRequest, app.Err("invalid token"))
	}
	tokenx, err := j.ParseWithClaims(t, &jwt.Claims{}, func(token *j.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil {
		return err
	}

	if !tokenx.Valid{
		return c.JSON(http.StatusBadRequest, app.Err("invalid token"))
	}

	token := jwt.From(tokenx)
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	go s.WsReader(conn, token)
	return nil
}


func (s *Server) AllRooms(c echo.Context) error {
	return c.JSON(http.StatusOK, s.Rooms)
}
