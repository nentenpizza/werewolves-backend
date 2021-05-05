package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type FriendsStorage struct {
	handler
}

func (s FriendsStorage) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	//g.GET("/", s.ByID)
	//Interaction with friends
	g.POST("/add", s.Add)
	//g.DELETE("/delete", s.Delete)

	//Interaction with friend requests
	g.POST("/accept", s.Accept)
	//g.POST("/decline", s.Decline)
	g.GET("/requests", s.Requests)

}

func (s FriendsStorage) ByID(c echo.Context) error {
	var form struct {
		UserID int `json:"user_id" validate:"required"`
	}
	err := c.Bind(&form)
	if err != nil {
		return err
	}
	err = c.Validate(&form)
	if err != nil {
		return err
	}

	friends, err := s.db.Friends.ByUserID(form.UserID, 0)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, friends)

}

func (s FriendsStorage) Add(c echo.Context) error {
	var form struct {
		TargetID int `json:"target_id" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return err
	}
	err = c.Validate(&form)
	if err != nil {
		return err
	}

	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	if f, _ := s.db.Friends.IsFriend(form.TargetID, user.ID); f {
		return c.JSON(http.StatusBadRequest, app.Err("you are already friends with this person"))
	}

	err = s.db.Friends.Add(storage.Friend{
		SenderID: user.ID,
		TargetID: form.TargetID,
		Active:   true,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}

//func (s FriendsStorage) Delete(c echo.Context) error {
//	var form struct {
//		TargetID int `json:"target_id" validate:"required"`
//	}
//
//	err := c.Bind(&form)
//	if err != nil {
//		return err
//	}
//	err = c.Validate(&form)
//	if err != nil {
//		return err
//	}
//
//	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)
//
//	if err != nil {
//		return err
//	}
//
//	if f, _ := s.db.Friends.IsFriend(form.TargetID, user.ID); !f {
//		return c.JSON(http.StatusBadRequest, app.Err("you are not a friend with this person"))
//	}
//
//	//s.db.Friends.Delete(storage.Friend{SenderID:})
//
//}

func (s FriendsStorage) Accept(c echo.Context) error {
	var form struct {
		SenderID int `json:"sender_id" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return err
	}
	err = c.Validate(&form)
	if err != nil {
		return err
	}

	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	if f, _ := s.db.Friends.IsFriend(form.SenderID, user.ID); f {
		return c.JSON(http.StatusBadRequest, app.Err("you are already friends with this person"))
	}

	err = s.db.Friends.Accept(storage.Friend{
		SenderID: form.SenderID,
		TargetID: user.ID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}

func (s FriendsStorage) Requests(c echo.Context) error {
	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	r, err := s.db.Friends.ByUserID(user.ID, 1)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)

}
