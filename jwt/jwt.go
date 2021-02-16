package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewWithClaims(claims Claims) *jwt.Token {
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func From(v interface{}) Claims {
	return *v.(*jwt.Token).Claims.(*Claims)
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"usr,omitempty"`
}