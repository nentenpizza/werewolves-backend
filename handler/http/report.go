package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
)

type ReportsEndpointGroup struct {
	handler
}

func (s ReportsEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/send", s.Report)
}

func (s ReportsEndpointGroup) Report(c echo.Context) error {
	var form struct {
		ReportedID int64  `json:"reported_id" validate:"required"`
		Reason     string `json:"reason" validate:"required,min=3"`
	}

	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}

	err := s.reportService.Report(jwt.From(c.Get("user")).ID, form.ReportedID, form.Reason)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}
