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
	DB            *storage.DB
	Secret        []byte
	AuthService   service.AuthService
	ReportService service.ReportService
	HonorService  service.HonorService
	FriendService service.FriendService
}

type handler struct {
	db            *storage.DB
	secret        []byte
	authService   service.AuthService
	reportService service.ReportService
	honorService  service.HonorService
	friendService service.FriendService
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:            h.DB,
		secret:        h.Secret,
		authService:   h.AuthService,
		reportService: h.ReportService,
		honorService:  h.HonorService,
		friendService: h.FriendService,
	}
}

func (h handler) Register(group *echo.Group, eg EndpointGroup) {
	eg.REGISTER(h, group)
}
