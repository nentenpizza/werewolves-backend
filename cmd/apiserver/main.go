package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/game/transport"
	"github.com/nentenpizza/werewolves/game/werewolves"
	"github.com/nentenpizza/werewolves/handler/http"
	"github.com/nentenpizza/werewolves/handler/websocket"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/validator"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
	"os"
	"time"
)

var uuid = []byte("d9799088-48bf-41c3-a109-6f09127f66bd")

var PGURL = flag.String("PG_URL", os.Getenv("PG_URL"), "url to your postgres db")

var phaseLength = flag.Int("phase", 30, "phase length in game")

func main() {
	flag.Parse()

	werewolves.PhaseLength = time.Duration(*phaseLength) * time.Second

	db, err := storage.Open(*PGURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	transport.Logger.SetOutput(os.Stdout)

	serv := service.New(db)

	wsHandler := transport.NewGame(
		transport.Game{
			DB: db,
			Clients: transport.NewClients(
				make(map[string]*transport.Client),
			),
			Rooms:          transport.NewRooms(make(map[string]*werewolves.Room)),
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
		http.GameEndpointGroup{PhaseLength: *phaseLength, Serv: server},
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
