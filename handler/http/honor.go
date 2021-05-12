package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
	"net/http"
)

type HonorsEndpointGroup struct {
	handler
}

func (s HonorsEndpointGroup) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/send", s.Honor)
}

func (s HonorsEndpointGroup) Honor(c echo.Context) error {
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

	err := s.reportService.Report(jwt.From(c.Get("user")).ID, form.HonoredID, form.Reason)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}
