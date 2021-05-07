package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/storage"
)

type Service interface {
	REGISTER(h handler, g *echo.Group)
}

type Handler struct {
	DB     *storage.DB
	Secret []byte
}

type handler struct {
	db *storage.DB
	s  []byte
}

func NewHandler(h Handler) *handler {
	return &handler{db: h.DB, s: h.Secret}
}

func (h handler) Register(group *echo.Group, service Service) {
	service.REGISTER(h, group)
}
