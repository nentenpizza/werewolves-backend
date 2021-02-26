package wserver

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type Context struct {
	Conn *websocket.Conn

	storage map[string]interface{}

	Update Update

	Token *jwt.Token
}

func (c *Context) Set(key string, val interface{}){
	c.storage[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.storage[key]
}

func (c *Context) Data() interface{}{
	return c.Update.Data
}

func (c *Context) EventType() string {
	return c.Update.EventType
}