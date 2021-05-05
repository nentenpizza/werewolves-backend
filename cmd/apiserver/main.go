package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/handler/http"
	"github.com/nentenpizza/werewolves/handler/websocket"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/validator"
	"github.com/nentenpizza/werewolves/werewolves"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
	"os"
	"time"
)

var uuid = []byte("d9799088-48bf-41c3-a109-6f09127f66bd")

var PGURL = flag.String("PG_URL", os.Getenv("PG_URL"), "url to your postgres db")

var phaseLength = flag.Int("phase", 5, "phase length in game")

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

	wsHandler := websocket.NewHandler(
		websocket.Handler{DB: db,
			Clients: websocket.NewClients(
				make(map[string]*websocket.Client),
			),
			Rooms:  websocket.NewRooms(make(map[string]*werewolves.Room)),
			Secret: uuid,
		},
	)

	server := wserver.NewServer(wserver.Settings{UseJWT: true, OnError: wsHandler.OnError, Claims: &jwt.Claims{}, Secret: uuid})

	server.Use(wsHandler.WebsocketJWT)
	server.Handle(websocket.EventTypeCreateRoom, wsHandler.OnCreateRoom)
	server.Handle(websocket.EventTypeJoinRoom, wsHandler.OnJoinRoom)
	server.Handle(websocket.EventTypeLeaveRoom, wsHandler.OnLeaveRoom)
	server.Handle(websocket.EventTypeStartGame, wsHandler.OnStartGame)
	server.Handle(wserver.OnConnect, wsHandler.OnConnect)
	server.Handle(websocket.EventTypeSendMessage, wsHandler.OnMessage)
	server.Handle(websocket.EventTypeVote, wsHandler.OnVote)
	server.Handle(websocket.EventTypeUseSkill, wsHandler.OnSkill)
	server.Handle(websocket.EventTypeSendEmote, wsHandler.OnEmote)

	h := http.NewHandler(
		http.Handler{
			DB: db,
		},
	)
	e := newEcho()
	//e.GET("/ws/:token", server.WsEndpoint)
	g := e.Group("", newJWTMiddleware())
	e.GET("/ws/:token", server.Listen)
	e.Static("/files", "assets")

	h.Register(
		e.Group("/api/auth"),
		http.AuthService{Secret: uuid})
	h.Register(
		g.Group("/api/users"),
		http.UsersService{},
	)
	h.Register(
		e.Group("/api/game"),
		http.GameService{},
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
