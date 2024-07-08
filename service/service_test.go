package service

import (
	"github.com/brianvoe/gofakeit/v6"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/nentenpizza/werewolves/storage"
	"github.com/pressly/goose"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
	"testing"
)

var db *storage.DB

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(7412))
	err := postgres.Start()
	if err != nil {
		postgres.Stop()
		panic(err)
	}

	db, err = storage.Open("host=localhost port=7412 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		postgres.Stop()
		panic(err)
	}

	if err := goose.Up(db.DB.DB, "../storage/migrate"); err != nil {
		postgres.Stop()
		panic(err)
	}

	code := m.Run()

	postgres.Stop()

	os.Exit(code)
}

func TestAuthService(t *testing.T) {
	var service AuthService

	Convey("Given auth service", t, func() {
		service = &Auth{
			Service{db: db},
		}

		form := SignUpForm{
			Email:    gofakeit.Email(),
			Username: "Username",
			Login:    "Login123",
			Password: gofakeit.Password(false, false, false, false, false, 16),
		}

		Convey("When the user entering valid data", func() {
			err := service.SignUp(form)
			So(err, ShouldEqual, nil)

		})

		Convey("When the user entering username with length greater than 10", func() {
			form.Username = strings.Repeat(gofakeit.Username(), 11)
			err := service.SignUp(form)
			So(err, ShouldEqual, InvalidUsername)

		})

		Convey("When the user entering username with length less than 3", func() {
			form.Username = "Us"
			err := service.SignUp(form)
			So(err, ShouldEqual, InvalidUsername)

		})

		Convey("When the user entering login with length greater than 16", func() {
			form.Login = strings.Repeat(gofakeit.Username(), 17)
			err := service.SignUp(form)
			So(err, ShouldEqual, InvalidLogin)

		})

		Convey("When the user entering login with length less than 3", func() {
			form.Login = "Us"
			err := service.SignUp(form)
			So(err, ShouldEqual, InvalidLogin)

		})
	})
}
