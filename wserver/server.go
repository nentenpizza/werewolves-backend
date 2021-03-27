package wserver

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/werewolves/app"
	j "github.com/nentenpizza/werewolves/jwt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	OnOther = iota
	OnConnect
	OnDisconnect
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 60,
}

var PongTimeout = time.Second * 60

type HandlerFunc func(ctx Context) error

type OnErrorFunc func(error, Context)

type Update struct {
	EventType string      `json:"event_type" mapstructure:"event_type"`
	Data      interface{} `json:"data" mapstructure:"data"`
}

type Settings struct {
	// secret for jwt, pass an empty string if u don't use it
	Secret []byte

	// if true, Context.Get("token") will return token
	UseJWT  bool
	OnError OnErrorFunc
	Claims  jwt.Claims
}

type Server struct {
	OnError  OnErrorFunc
	handlers map[interface{}]HandlerFunc
	useJWT   bool
	secret   []byte
	claims   jwt.Claims
}

type Conn struct {
	conn *websocket.Conn
	sync.Mutex
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c *Conn) WriteMessage(code int, msg []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.conn.WriteMessage(code, msg)
}

func (c *Conn) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Conn) SetPongHandler(h func(appdata string) error) {
	c.Lock()
	defer c.Unlock()
	c.conn.SetPongHandler(h)
}

func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	return c.Close()
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	c.Lock()
	defer c.Unlock()
	return c.ReadMessage()
}

func NewServer(s Settings) *Server {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	if s.UseJWT && len(s.Secret) < 1 {
		panic("wserver: secret can not be empty string if UseJWT enabled")
	}
	return &Server{
		useJWT:   s.UseJWT,
		OnError:  s.OnError,
		handlers: make(map[interface{}]HandlerFunc),
		secret:   s.Secret,
		claims:   s.Claims,
	}
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func (s *Server) Handle(eventType interface{}, h HandlerFunc, middleware ...MiddlewareFunc) {
	switch eventType.(type) {
	case int:
		break
	case string:
		break
	default:
		panic("wserver: unsupported eventType")
	}
	h = applyMiddleware(h, middleware...)
	s.handlers[eventType] = h
}

func (s *Server) runHandler(h HandlerFunc, c Context) {
	f := func() {
		if err := h(c); err != nil {
			if s.OnError != nil {
				s.OnError(err, c)
			} else {
				log.Println(err)
			}
		}
	}
	f()
}

// Listen is handler that upgrades http client to websocket client
func (s *Server) Listen(c echo.Context) error {
	var tok string
	var err error
	if s.useJWT {
		tok = c.Param("token")
		if tok == "" {
			return c.JSON(http.StatusBadRequest, app.Err("invalid token"))
		}
		tokenx, err := jwt.ParseWithClaims(tok, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, "bad token")
		}
		if !tokenx.Valid {
			return c.JSON(http.StatusBadRequest, "bad token")
		}

	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	conn := NewConn(ws)
	go s.keepAlive(conn, PongTimeout)
	go s.reader(conn, tok)
	return nil
}

func (s *Server) keepAlive(conn *Conn, timeout time.Duration) {
	lastResponse := time.Now()
	conn.SetPongHandler(func(_ string) error {
		lastResponse = time.Now()
		return nil
	})
	for {
		err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
		if err != nil {
			return
		}
		time.Sleep((timeout * 9) / 10)
		if time.Now().Sub(lastResponse) > timeout {
			log.Printf("Ping don't get response, disconnecting to %s", conn.conn.LocalAddr())
			err = conn.Close()
			if s.OnError != nil {
				s.OnError(err, Context{})
			}
			return
		}
	}
}

func (s *Server) reader(conn *Conn, token string) {
	ctx := Context{Conn: conn, storage: make(map[string]interface{})}
	ctx.Set("token", token)
	s.runOnConnectHandler(ctx)
	for {
		ctx := Context{Conn: conn, storage: make(map[string]interface{})}
		ctx.Set("token", token)
		_, msg, err := conn.conn.ReadMessage()
		if err != nil {
			s.OnError(err, ctx)
			s.runOnDisconnectHandler(ctx)
			return
		}

		s.processUpdate(msg, ctx)
	}
}

func (s Server) runOnDisconnectHandler(ctx Context) {
	h, ok := s.handlers[OnDisconnect]
	if ok {
		s.runHandler(h, ctx)
	}
}

func (s *Server) runOnConnectHandler(ctx Context) {
	h, ok := s.handlers[OnConnect]
	if ok {
		s.runHandler(h, ctx)
	}
}

func (s *Server) processUpdate(msg []byte, c Context) {
	u := &Update{}
	err := json.Unmarshal(msg, u)
	if err != nil {
		if s.OnError != nil {
			s.OnError(err, c)
		}
	}
	c.Update = u

	handler, ok := s.handlers[u.EventType]
	if !ok {
		h, ok := s.handlers[OnOther]
		if ok {
			handler = h
		} else {
			return
		}
	}
	if handler != nil {
		s.runHandler(handler, c)
	}
}
