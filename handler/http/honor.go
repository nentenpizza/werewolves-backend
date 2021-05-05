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
	g.Use()
	g.POST("/send", s.Honor)
}

func (s HonorsService) Honor(c echo.Context) error {
	var form struct {
		HonoredID int    `json:"honored_id" validate:"required"`
		Reason    string `json:"reason" validate:"required"`
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

	if exists, _ := s.db.Users.ExistsByID(form.HonoredID); !exists {
		return c.JSON(http.StatusBadRequest, app.Err("user does not exist"))
	}

	if user.ID == form.HonoredID {
		return c.JSON(http.StatusBadRequest, app.Err("you cannot honor myself"))
	}

	if exists, _ := s.db.Honors.Exists(form.HonoredID, user.ID); exists {
		return c.JSON(http.StatusForbidden, app.Err("you already honored this user"))
	}

	if count, _ := s.db.Honors.CountToday(user.ID); count > 10 {
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
