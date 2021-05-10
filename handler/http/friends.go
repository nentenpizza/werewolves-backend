package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"net/http"
)

type FriendsService struct {
	handler
}

func (s FriendsService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/request", s.Request)
	g.POST("/accept", s.Accept)
	g.GET("/list", s.Friends)
}

func (s FriendsService) Request(c echo.Context) error {
	var form struct {
		Receiver int64 `json:"id"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	claims := jwt.From(c.Get("user"))
	if form.Receiver == claims.ID {
		return c.JSON(http.StatusBadRequest, app.Err("you cannot request yourself"))
	}
	me, err := s.db.Users.ByID(claims.ID)
	if err != nil {
		return err
	}
	has, err := s.db.Friends.IsFriend(me.Relations, form.Receiver)
	if err != nil {
		return err
	}
	if has {
		return c.JSON(http.StatusConflict, app.Err("receiver already your friend"))
	}
	id, err := s.db.Friends.Create(claims.ID)
	if err != nil {
		return err
	}
	err = s.db.Users.UpdateRelations(form.Receiver, id)
	if err != nil {
		return err
	}
	err = s.db.Users.UpdateRelations(claims.ID, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"id": id})
}

func (s FriendsService) Accept(c echo.Context) error {
	var form struct {
		ID int `json:"id"`
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

func (s FriendsService) Friends(c echo.Context) error {
	claims := jwt.From(c.Get("user"))
	user, err := s.db.Users.ByID(claims.ID)
	if err != nil {
		return err
	}
	users, err := s.db.Friends.UsersByID(user.Relations, user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}
