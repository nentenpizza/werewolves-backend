package service

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/nentenpizza/werewolves/app"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(form SignUpForm) error
	SignIn(form SignInForm, secret []byte) (string, error)
}

type Auth struct {
	Service
}

type SignUpForm struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s Auth) SignUp(form SignUpForm) error {
	form.Login = strings.ToLower(form.Login)
	exists, err := s.db.Users.Exists(form.Username)
	if err != nil {
		return err
	}
	if exists {
		return serviceError(http.StatusConflict, "username already taken")
	}

	exists, err = s.db.Users.ExistsByLogin(form.Login)
	if err != nil {
		return err
	}
	if exists {
		return serviceError(http.StatusConflict, "login already taken")
	}
	if !s.validateUsername(form.Username) {
		return InvalidUsername
	}
	if !s.validateLogin(form.Login) {
		return InvalidLogin
	}
	if !s.validatePassword(form.Password) {
		return InvalidPassword
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
	return err
}

func (s Auth) SignIn(form SignInForm, secret []byte) (string, error) {
	form.Login = strings.ToLower(form.Login)
	user, err := s.db.Users.ByLogin(form.Login)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		return "", serviceError(http.StatusNotFound, "user not found")
	}

	if !s.compareHash(form.Password, user.EncryptedPassword) {
		return "", serviceError(http.StatusBadRequest, "wrong credentials")
	}

	if user.BannedUntil.Sub(time.Now()).Seconds() > 0 {
		return "", serviceError(http.StatusForbidden, "user is restricted")
	}

	token := jwt.NewWithClaims(jwt.Claims{
		Username: user.Username,
		ID:       user.ID,
	})

	t, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (s Auth) compareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s Auth) validateUsername(username string) bool {
	if len(username) > 10 || len(username) <= 3 {
		return false
	}
	if app.StringContains(username, "lives", "matter") {
		return false
	}
	return true
}

func (s Auth) validateLogin(login string) bool {
	if len(login) > 16 || len(login) <= 3 {
		return false
	}
	return true
}

func (s Auth) validatePassword(password string) bool {
	if len(password) < 5 && len(password) > 30 {
		return false
	}
	return true
}
