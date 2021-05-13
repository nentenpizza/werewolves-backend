package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/handler/http"
	"github.com/nentenpizza/werewolves/handler/websocket"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/service"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/validator"
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"io"
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

	logFile, err := os.Create(fmt.Sprintf("logs_%d", time.Now().Unix()))
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	websocket.Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	wsHandler := websocket.NewHandler(
		websocket.Handler{
			DB: db,
			Clients: websocket.NewClients(
				make(map[string]*websocket.Client),
			),
			Rooms:  websocket.NewRooms(make(map[string]*werewolves.Room)),
			Secret: uuid,
		},
	)

	server := wserver.NewServer(wserver.Settings{UseJWT: true, OnError: wsHandler.OnError, Claims: &jwt.Claims{}, Secret: uuid})

	server.Use(wsHandler.WebsocketJWT, wsHandler.Logger)
	server.Handle(websocket.EventTypeCreateRoom, wsHandler.OnCreateRoom)
	server.Handle(websocket.EventTypeJoinRoom, wsHandler.OnJoinRoom)
	server.Handle(websocket.EventTypeLeaveRoom, wsHandler.OnLeaveRoom)
	server.Handle(websocket.EventTypeStartGame, wsHandler.OnStartGame)

	server.Handle(wserver.OnConnect, wsHandler.OnConnect)
	server.Handle(wserver.OnDisconnect, wsHandler.OnDisconnect)

	server.Handle(websocket.EventTypeSendMessage, wsHandler.OnMessage)
	server.Handle(websocket.EventTypeVote, wsHandler.OnVote)
	server.Handle(websocket.EventTypeUseSkill, wsHandler.OnSkill)
	server.Handle(websocket.EventTypeSendEmote, wsHandler.OnEmote)

	serv := service.NewService(db)
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
	//e.GET("/ws/:token", server.WsEndpoint)
	g := e.Group("", newJWTMiddleware())
	e.Static("/files", "assets")

	h.Register(
		e.Group("/api/auth"),
		http.AuthEndpointGroup{})
	h.Register(
		g.Group("/api/users"),
		http.UsersEndpointGroup{},
	)
	h.Register(
		e.Group(""),
		http.GameEndpointGroup{PhaseLength: *phaseLength, Serv: server},
	)

	h.Register(
		g.Group("/api/reports"),
		http.ReportsEndpointGroup{},
	)

	h.Register(
		g.Group("/api/honors"),
		http.HonorsEndpointGroup{},
	)

	h.Register(
		g.Group("/api/inventory"),
		http.ItemsEndpointGroup{},
	)

	h.Register(
		g.Group("/api/friends"),
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
