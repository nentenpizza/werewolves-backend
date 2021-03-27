package http

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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
	g.GET("/check", s.CheckToken)
}

// Register is endpoint for signing in
func (s AuthService) Register(c echo.Context) error {
	var form struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	form.Login = strings.ToLower(form.Login)
	exists, err := s.db.Users.Exists(form.Username)
	if err != nil {
		return err
	}
	if exists {
		return c.JSON(http.StatusConflict, app.Err("username already taken"))
	}

	exists, err = s.db.Users.ExistsByLogin(form.Username)
	if err != nil {
		return err
	}
	if exists {
		return c.JSON(http.StatusConflict, app.Err("login already taken"))
	}
	if !s.validateUsername(form.Username) {
		return c.JSON(http.StatusBadRequest, app.Err("usr 3-10 chars"))
	}
	if !s.validateLogin(form.Login) {
		return c.JSON(http.StatusBadRequest, app.Err("login 3-16 chars"))
	}
	if !s.validatePassword(form.Password) {
		return c.JSON(http.StatusBadRequest, app.Err("password 5-30 chars"))
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u := storage.User{
		Email:             form.Email,
		Username:          form.Username,
		Login:             form.Login,
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
	if len(username) > 10 || len(username) <= 3 {
		return false
	}
	if app.StringContains(username, "lives", "matter") {
		return false
	}
	return true
}

func (s AuthService) validateLogin(login string) bool {
	if len(login) > 16 || len(login) <= 3 {
		return false
	}
	return true
}

func (s AuthService) validatePassword(password string) bool {
	if len(password) < 5 && len(password) > 30 {
		return false
	}
	return true
}

// Login is endpoint for logging in
// Not done yet.
func (s AuthService) Login(c echo.Context) error {
	var form struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	form.Login = strings.ToLower(form.Login)
	user, err := s.db.Users.ByLogin(form.Login)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return c.JSON(http.StatusNotFound, app.Err("user not found"))
	}

	if !s.compareHash(form.Password, user.EncryptedPassword) {
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

func (s AuthService) CheckToken(c echo.Context) error {
	return c.JSON(200, app.Ok())
}

func (s AuthService) compareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
