package websocket

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	j "github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/wserver"
	"log"
)

func (h *handler) WebsocketJWT() wserver.MiddlewareFunc {
	return func(next wserver.HandlerFunc) wserver.HandlerFunc {
		return func(c *wserver.Context) error {
			tok := (c.Get("token")).(string)
			if tok == "" {
				return errors.New("token is nil")
			}
			tokenx, err := jwt.ParseWithClaims(tok, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return h.s, nil
			})
			if err != nil {
				return err
			}

			if !tokenx.Valid {
				return errors.New("middleware: invalid token")
			}
			token := j.From(tokenx)
			client := h.c.Read(token.Username)
			log.Println(token.Username, "Token")
			if client != nil {
				client.conn = c.Conn
				if client.Token.Username == "" {
					client.Token = token
				}
			} else {
				client = NewClient(c.Conn, token, make([]interface{}, 0), make(chan bool))
				h.c.Write(client.Token.Username, client)
			}
			c.Set("client", client)
			return next(c)
		}
	}
}
