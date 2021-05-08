package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type ReportsService struct {
	handler
}

func (s ReportsService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/send", s.Report)
}

func (s ReportsService) Report(c echo.Context) error {
	var form struct {
		ReportedID int    `json:"reported_id" validate:"required"`
		Reason     string `json:"reason" validate:"required,min=3"`
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

	exists, err := s.db.Users.ExistsByID(form.ReportedID)
	if err != nil {
		return err
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, app.Err("user does not exist"))
	}

	if form.ReportedID == user.ID {
		return c.JSON(http.StatusBadRequest, app.Err("you cannot report yourself"))
	}

	err = s.db.Reports.Create(storage.Report{
		ReportedID: form.ReportedID,
		Reason:     form.Reason,
		SenderID:   user.ID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, app.Ok())

}
