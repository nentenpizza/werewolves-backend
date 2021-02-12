package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/storage"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
}

type AuthService struct {
	handler
	Secret []byte
}

func (s AuthService) REGISTER(h handler, g *echo.Group) {
	s.handler = h
	g.POST("/login", s.Login)
	g.POST("/register", s.Login)
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
	return err
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
	return nil
}

// CompareHash returns nil on success
func (s AuthService) CompareHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}