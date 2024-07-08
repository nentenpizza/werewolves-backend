package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
)

type ItemsEndpointGroup struct {
	handler
}

func (s ItemsEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h
	g.GET("/items", s.Items)
	g.POST("/count", s.Count)
	g.PUT("/item", s.Item)
	g.DELETE("/item", s.ItemDelete)

}

func (s ItemsEndpointGroup) Items(c echo.Context) error {
	items, err := s.db.Items.Items(jwt.From(c.Get("user")).ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"items": items})

}

func (s ItemsEndpointGroup) Item(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}

	count, err := s.itemService.GiveItem(jwt.From(c.Get("user")).ID, form.Item)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})

}

func (s ItemsEndpointGroup) ItemDelete(c echo.Context) error {
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

	count, err := s.itemService.DeleteItem(jwt.From(c.Get("user")).ID, form.Item, form.Count)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"count": count})

}

func (s ItemsEndpointGroup) Count(c echo.Context) error {
	var form struct {
		Item string `json:"item" validate:"required"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}

	count, err := s.itemService.Count(jwt.From(c.Get("user")).ID, form.Item)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})

}
