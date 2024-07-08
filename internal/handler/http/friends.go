package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/internal/service"
	"github.com/nentenpizza/werewolves/pkg/jwt"
)

type FriendsEndpointGroup struct {
	handler
}

func (s FriendsEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h

	g.POST("/request", s.Request)
	g.POST("/accept", s.Accept)
	g.POST("/accept_userid", s.AcceptByUserID)

	g.GET("/list", s.Friends)
	g.GET("/list_unaccepted", s.ListUnaccepted)
}

func (s FriendsEndpointGroup) Request(c echo.Context) error {
	var form struct {
		Receiver int64 `json:"id"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	claims := jwt.From(c.Get("user"))
	id, err := s.friendService.Request(claims.ID, form.Receiver)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"id": id})
}

func (s FriendsEndpointGroup) Accept(c echo.Context) error {
	var form struct {
		ID int `json:"request_id"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	claims := jwt.From(c.Get("user"))
	err := s.db.Friends.Accept(claims.ID, form.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, app.Ok())
}

func (s FriendsEndpointGroup) Friends(c echo.Context) error {
	users, err := s.friendService.UserFriends(jwt.From(c.Get("user")).ID)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"friends": users})
}

func (s FriendsEndpointGroup) ListUnaccepted(c echo.Context) error {
	users, err := s.friendService.UnacceptedUsers(jwt.From(c.Get("user")).ID)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"unaccepted": users})
}

func (s FriendsEndpointGroup) AcceptByUserID(c echo.Context) error {
	var form struct {
		ID int64 `json:"user_id"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	err := s.friendService.AcceptBySenderID(form.ID, jwt.From(c.Get("user")).ID)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, app.Ok())
}
