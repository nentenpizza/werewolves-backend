package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/storage"
)

type Service interface{
	REGISTER(h handler, g *echo.Group)
}

type Handler struct {
	DB *storage.DB
}

type handler struct {
	db *storage.DB
}

func NewHandler(h Handler) *handler {
	return &handler{db: h.DB}
}

func (h handler) Register(group *echo.Group, service Service)  {
	service.REGISTER(h, group)
}