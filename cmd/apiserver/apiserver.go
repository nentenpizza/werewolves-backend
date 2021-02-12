package apiserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/werewolves/handler"
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
	h := handler.New(
		handler.Handler{
		DB:db,
		})
	e := newEcho()
	h.Register(
		e.Group("/api/auth"),
		handler.AuthService{Secret: uuid})
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