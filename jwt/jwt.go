package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

func NewWithClaims(claims Claims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func From(v interface{}) Claims {
	return *v.(*jwt.Token).Claims.(*Claims)
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"usr,omitempty"`
}