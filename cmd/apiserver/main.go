package main

import (
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
)

var uuid = []byte("d9799088-48bf-41c3-a109-6f09127f66bd")


func onError(err error, c wserver.Context){
	log.Println(err, c)
}

func middle() wserver.MiddlewareFunc{
	return func(next wserver.HandlerFunc) wserver.HandlerFunc{
		return func(c wserver.Context) error {
			log.Println("incoming message")
			return next(c)
		}
	}
}

func main(){
	db, err := storage.Open(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err !=nil{
		log.Fatal(err)
	}
	defer db.Close()

	wsHandler := websocket.NewHandler(
		websocket.Handler{DB: db,
			Clients: websocket.NewClients(
			make(map[string]*websocket.Client),
			),
			Rooms: websocket.NewRooms(make(map[string]*werewolves.Room)),
			Secret: uuid,
		},
		)

	server := wserver.NewServer(wserver.Settings{UseJWT: true, OnError: onError, Claims: &jwt.Claims{}, Secret: uuid})
	server.Handle(websocket.EventTypeCreateRoom, wsHandler.OnCreateRoom, wsHandler.WebsocketJWT())
	server.Handle(websocket.EventTypeJoinRoom, wsHandler.OnJoinRoom, wsHandler.WebsocketJWT())
	server.Handle(websocket.EventTypeLeaveRoom, wsHandler.OnLeaveRoom, wsHandler.WebsocketJWT())
	server.Handle(websocket.EventTypeStartGame, wsHandler.OnStartGame, wsHandler.WebsocketJWT())
	server.Handle(wserver.OnConnect, wsHandler.OnConnect, wsHandler.WebsocketJWT())
	server.Handle(websocket.EventTypeSendMessage, wsHandler.OnMessage, wsHandler.WebsocketJWT())

	h := http.NewHandler(
		http.Handler{
			DB:db,
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
	e.Logger.Fatal(e.Start(":7070"))
}


func newEcho() *echo.Echo{
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