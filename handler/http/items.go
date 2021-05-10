package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type ItemsService struct {
	handler
}

func (s ItemsService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/items", s.Items)
	g.POST("/count", s.Count)
	g.PUT("/item", s.Item)
	g.DELETE("/item", s.ItemDelete)

}

func (s ItemsService) Items(c echo.Context) error {
	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	inv, err := s.db.Items.Items(user.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, inv)

}

func (s ItemsService) Item(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
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

	inv := storage.Item{
		Name:   form.Item,
		UserID: user.ID,
	}

	count, err := s.db.Items.Count(inv)

	if err != nil {
		return err
	}

	err = s.db.Items.Create(inv)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"ok": count + 1})

}

func (s ItemsService) ItemDelete(c echo.Context) error {
	var form struct {
		Item  string `json:"item" validate:"required"`
		Count uint   `json:"count" validate:"required"`
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

	inv := storage.Item{
		UserID: user.ID,
		Name:   form.Item,
	}

	err = s.db.Items.Delete(inv, form.Count)
	if err != nil {
		return err
	}

	count, err := s.db.Items.Count(inv)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"ok": count})

}

func (s ItemsService) Count(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
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

	count, err := s.db.Items.Count(storage.Item{
		Name:   form.Item,
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"ok": count})

}
