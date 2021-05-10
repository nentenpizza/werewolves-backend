package websocket

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	j "github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/wserver"
	log "github.com/sirupsen/logrus"
)

func (h *handler) WebsocketJWT(next wserver.HandlerFunc) wserver.HandlerFunc {
	return func(c *wserver.Context) error {
		tok := (c.Get("token")).(*jwt.Token)
		if tok == nil {
			return errors.New("token is nil")
		}
		token := j.From(tok)
		client := h.c.Read(token.Username)
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

func (h *handler) Logger(next wserver.HandlerFunc) wserver.HandlerFunc {
	return func(c *wserver.Context) error {
		client, ok := c.Get("client").(*Client)
		if ok {
			if c.Update != nil {
				Logger.WithFields(log.Fields{
					"event_type":    c.Update.EventType,
					"event":         c.Update.Data,
					"from_username": client.Token.Username,
				}).Info("incoming message")
			} else {
				Logger.WithFields(log.Fields{
					"from_username": client.Token.Username,
				}).Info("connecting")
			}
		} else {
			Logger.WithFields(log.Fields{
				"event_type": c.Update.EventType,
				"event":      c.Update.Data,
			}).Info("incoming message")
		}
		return next(c)
	}
}
