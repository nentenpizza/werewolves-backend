package wserver

import (
	"github.com/mitchellh/mapstructure"
)

type Context struct {
	Conn *Conn

	storage map[string]interface{}

	Update *Update
}

func NewContext(conn *Conn) *Context {
	return &Context{Conn: conn, storage: make(map[string]interface{})}
}

func (c *Context) Set(key string, val interface{}) {
	c.storage[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.storage[key]
}

func (c *Context) Bind(i interface{}) error {
	err := mapstructure.Decode(c.Update.Data, i)
	return err
}
func (c *Context) Data() interface{} {
	return c.Update.Data
}

func (c *Context) EventType() string {
	return c.Update.EventType
}
