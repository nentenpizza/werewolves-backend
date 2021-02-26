package websocket

import (
	"errors"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/wserver"
)

func (s *handler) WebsocketJWT() wserver.MiddlewareFunc {
	return func (next wserver.HandlerFunc) wserver.HandlerFunc {
		return func(c wserver.Context) error {
			tok := c.Token
			if tok == nil{
				return errors.New("token is nil")
			}
			token := jwt.From(tok)
			if token.Username != "" {
				client := s.c.Read(token.Username)
				if client != nil {
					client.Token = token
				}else{
					client = NewClient(c.Conn, token, make([]interface{}, 0), make(chan bool))
				}
				c.Set("client", client)
			} else {
				return errors.New("token invalid")
			}
			return next(c)
		}
	}
}

