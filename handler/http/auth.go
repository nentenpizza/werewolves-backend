package http

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthService struct {
	handler
	Secret []byte
}

func (s AuthService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/login", s.Login)
	g.POST("/register", s.Register)
}

// Register is endpoint for signing in
func (s AuthService) Register(c echo.Context) error {
	var form struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	exists, err := s.db.Users.Exists(form.Username)
	if err != nil {
		return err
	}
	if exists{
		return c.JSON(http.StatusConflict, app.Err("username already taken"))
	}
	if !s.validateUsername(form.Username){
		return c.JSON(http.StatusBadRequest, app.Err("username must contains less than 16 symbols"))
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u := storage.User{
		Email: form.Email,
		Username: form.Username,
		EncryptedPassword: string(encryptedPass),
	}
	err = s.db.Users.Create(u)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, u)
	return err
}

func (s AuthService) validateUsername(username string) bool {
	if len(username ) > 15 {return false}
	if app.StringContains(username, "lives", "matter"){
		return false
	}
	return true
}


// Login is endpoint for logging in
// Not done yet.
func (s AuthService) Login(c echo.Context) error {
	var form struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	user, err := s.db.Users.ByUsername(form.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return c.JSON(http.StatusNotFound, app.Err("user not found"))
	}

	if !s.compareHash(form.Password,user.EncryptedPassword) {
		return c.JSON(http.StatusBadRequest, app.Err("wrong credentials"))
	}

	if user.BannedUntil.Sub(time.Now()).Seconds() > 0 {
		return c.JSON(http.StatusForbidden, app.Err("user is restricted"))
	}

	token := jwt.NewWithClaims(jwt.Claims{
		Username: user.Username,
	})

	t, err := token.SignedString(s.Secret)
	if err != nil {
		return err
	}

	return c.JSON(200, echo.Map{"token": t})
}

func (s AuthService) compareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err != nil{
		return false
	}
	return true
}