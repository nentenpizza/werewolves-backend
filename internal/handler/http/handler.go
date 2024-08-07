package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/internal/service"
	"github.com/nentenpizza/werewolves/internal/storage"
)

type EndpointGroup interface {
	Register(h handler, g *echo.Group)
}

type Handler struct {
	DB            *storage.DB
	Secret        []byte
	AuthService   service.AuthService
	ReportService service.ReportService
	HonorService  service.HonorService
	FriendService service.FriendService
	ItemService   service.ItemService
}

type handler struct {
	db            *storage.DB
	secret        []byte
	authService   service.AuthService
	reportService service.ReportService
	honorService  service.HonorService
	friendService service.FriendService
	itemService   service.ItemService
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:            h.DB,
		secret:        h.Secret,
		authService:   h.AuthService,
		reportService: h.ReportService,
		honorService:  h.HonorService,
		friendService: h.FriendService,
		itemService:   h.ItemService,
	}
}

func (h handler) Register(group *echo.Group, eg EndpointGroup) {
	eg.Register(h, group)
}
