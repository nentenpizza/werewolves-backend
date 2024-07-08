package main

import (
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/handler/http"
	"github.com/nentenpizza/werewolves/handler/websocket"
	"github.com/nentenpizza/werewolves/internal/werewolves"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/validator"
	"github.com/nentenpizza/werewolves/wserver"
)

var uuid = []byte("d9799088-48bf-41c3-a109-6f09127f66bd") // for dev purposes (i forgot xDDD)

func main() {

	phaseLength, err := strconv.Atoi(os.Getenv("PHASE_LENGTH"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.Open(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	werewolves.Logger.SetOutput(os.Stdout)

	serv := service.New(db)

	wsHandler := werewolves.NewGame(
		werewolves.Game{
			DB: db,
			Clients: werewolves.NewClients(
				make(map[string]*werewolves.Client),
			),
			Rooms:          werewolves.NewRooms(),
			Secret:         uuid,
			FriendsService: &service.Friends{Service: serv},
		},
	)

	server := websocket.Initialize(wserver.Settings{
		UseJWT: true, OnError: wsHandler.OnError, Claims: &jwt.Claims{}, Secret: uuid,
	},
		wsHandler,
	)

	h := http.NewHandler(
		http.Handler{
			DB:            db,
			Secret:        uuid,
			AuthService:   &service.Auth{Service: serv},
			ReportService: &service.Reports{Service: serv},
			HonorService:  &service.Honors{Service: serv},
			FriendService: &service.Friends{Service: serv},
			ItemService:   &service.Items{Service: serv},
		},
	)
	e := newEcho()

	e.Static("/files", "assets")

	api := e.Group("/api", newJWTMiddleware())

	h.Register(
		e.Group("/api/auth"),
		http.AuthEndpointGroup{})
	h.Register(
		api.Group("/users"),
		http.UsersEndpointGroup{},
	)
	h.Register(
		e.Group(""),
		http.GameEndpointGroup{PhaseLength: phaseLength, Serv: server},
	)

	h.Register(
		api.Group("/reports"),
		http.ReportsEndpointGroup{},
	)

	h.Register(
		api.Group("/honors"),
		http.HonorsEndpointGroup{},
	)

	h.Register(
		api.Group("/inventory"),
		http.ItemsEndpointGroup{},
	)

	h.Register(
		api.Group("/friends"),
		http.FriendsEndpointGroup{},
	)

	e.Logger.Fatal(e.Start(":7070"))
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Validator = validator.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Println(err)
		e.DefaultHTTPErrorHandler(err, c)
	}
	return e
}

func newJWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: uuid,
		Claims:     &jwt.Claims{},
	})
}
