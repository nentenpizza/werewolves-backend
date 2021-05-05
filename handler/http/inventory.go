package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type InventoryService struct {
	handler
}

func (s InventoryService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/fetch", s.Fetch)
	g.POST("/count", s.CountItem)
	g.PUT("/item", s.Item)
	g.DELETE("/item", s.ItemDelete)

}

func (s InventoryService) Fetch(c echo.Context) error {
	user, err := s.db.Users.ByUsername(jwt.From(c.Get("user")).Username)

	if err != nil {
		return err
	}

	inv, err := s.db.Inventory.Fetch(user.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, inv)

}

func (s InventoryService) Item(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
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

	inv := storage.Inventory{
		Item:   form.Item,
		UserID: user.ID,
	}

	count, err := s.db.Inventory.CountItem(inv)

	if err != nil {
		return err
	}

	err = s.db.Inventory.Create(inv)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.OkWithField(count+1))

}

func (s InventoryService) ItemDelete(c echo.Context) error {
	var form struct {
		Item  string `json:"item" validate:"required"`
		Count uint   `json:"count" validate:"required"`
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

	inv := storage.Inventory{
		UserID: user.ID,
		Item:   form.Item,
	}

	err = s.db.Inventory.Delete(inv, form.Count)
	if err != nil {
		return err
	}

	count, err := s.db.Inventory.CountItem(inv)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.OkWithField(count))

}

func (s InventoryService) CountItem(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
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

	count, err := s.db.Inventory.CountItem(storage.Inventory{
		Item:   form.Item,
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.OkWithField(count))

}
