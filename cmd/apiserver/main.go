package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/handler"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/nentenpizza/werewolves/validator"
	"log"
	"os"
)

var uuid = []byte("d9799088-48bf-41c3-a109-6f09127f66bd")


func main(){
	db, err := storage.Open(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err !=nil{
		log.Fatal(err)
	}
	defer db.Close()
	server := handler.NewServer(uuid)
	h := handler.New(
		handler.Handler{
			DB:db,
		},
		)
	e := newEcho()
	e.GET("/ws/:token", server.WsEndpoint)
	g := e.Group("", newJWTMiddleware())


	h.Register(
		e.Group("/api/auth"),
		handler.AuthService{Secret: uuid})
	h.Register(
		g.Group("/api"),
		server,
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