package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/service"
)

type AuthEndpointGroup struct {
	handler
}

func (s AuthEndpointGroup) Register(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/login", s.Login)
	g.POST("/RegisterUser", s.RegisterUser)
	g.GET("/check", s.CheckToken)
}

// Register is endpoint for signing in
func (s AuthEndpointGroup) RegisterUser(c echo.Context) error {
	var form service.SignUpForm
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	err := s.authService.SignUp(form)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}
	return c.JSON(http.StatusCreated, app.Ok())
}

// Login is endpoint for logging in
func (s AuthEndpointGroup) Login(c echo.Context) error {
	var form service.SignInForm
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	t, err := s.authService.SignIn(form, s.secret)
	if err != nil {
		serviceErr, ok := err.(*service.Error)
		if ok {
			return c.JSON(serviceErr.Code, serviceErr.Error())
		}
		return err
	}

	return c.JSON(200, echo.Map{"token": t})
}

func (s AuthEndpointGroup) CheckToken(c echo.Context) error {
	return c.JSON(200, app.Ok())
}
