package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/service"
	"github.com/nentenpizza/werewolves/storage"
)

type EndpointGroup interface {
	REGISTER(h handler, g *echo.Group)
}

type Handler struct {
	DB          *storage.DB
	AuthService service.AuthService
	Secret      []byte
}

type handler struct {
	db          *storage.DB
	authService service.AuthService
	s           []byte
}

func NewHandler(h Handler) *handler {
	return &handler{db: h.DB, authService: h.AuthService, s: h.Secret}
}

func (h handler) Register(group *echo.Group, eg EndpointGroup) {
	eg.REGISTER(h, group)
}
