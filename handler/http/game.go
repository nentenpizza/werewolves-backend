package http

import (
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	j "github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/wserver"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 60,
}

type GameEndpointGroup struct {
	handler
	PhaseLength int
	Serv        *wserver.Server
}

func (s GameEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	api := g.Group("/api/game")

	api.GET("/phase_length", s.GetPhaseLength)

	g.GET("/ws/:token", s.GameEndpoint)
}

func (s GameEndpointGroup) GetPhaseLength(c echo.Context) error {
	return c.JSON(200, echo.Map{"length": s.PhaseLength})
}

func (s *GameEndpointGroup) GameEndpoint(c echo.Context) error {
	var token *jwtgo.Token
	var err error
	tok := c.Param("token")
	if tok == "" {
		return c.JSON(http.StatusBadRequest, app.Err("invalid token"))
	}
	token, err = jwtgo.ParseWithClaims(tok, &j.Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad token")
	}
	if !token.Valid {
		return c.JSON(http.StatusBadRequest, "bad token")
	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	s.Serv.Accept(ws, token)
	return nil
}
