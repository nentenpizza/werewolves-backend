package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type HonorsService struct {
	Secret []byte
	handler
}

func (s HonorsService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/send", s.Honor)
}

func (s HonorsService) Honor(c echo.Context) error {
	var form struct {
		HonoredID int64  `json:"honored_id" validate:"required"`
		Reason    string `json:"reason" validate:"required"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}

	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	exists, err := s.db.Users.ExistsByID(form.HonoredID)
	if err != nil {
		return err
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, app.Err("user does not exist"))
	}

	if user.ID == form.HonoredID {
		return c.JSON(http.StatusBadRequest, app.Err("you cannot honor yourself"))
	}

	exists, err = s.db.Honors.Exists(form.HonoredID, user.ID)
	if err != nil {
		return err
	}
	if exists {
		return c.JSON(http.StatusForbidden, app.Err("you already honored this user"))
	}

	count, err := s.db.Honors.CountToday(user.ID)
	if err != nil {
		return err
	}
	if count > 10 {
		return c.JSON(http.StatusForbidden, app.Err("you reached daily limit"))
	}

	err = s.db.Honors.Create(storage.Honor{
		HonoredID: form.HonoredID,
		Reason:    form.Reason,
		SenderID:  user.ID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}
