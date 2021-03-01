package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/nentenpizza/werewolves/jwt"
	"github.com/nentenpizza/werewolves/werewolves"
	"log"
	"sync"
)

type Client struct {
	sync.Mutex
	conn *websocket.Conn
	*werewolves.Player
	room *werewolves.Room
	AFK bool
	Token jwt.Claims
	Unreached []interface{}
	quit chan bool
}

func NewClient(conn *websocket.Conn, token jwt.Claims, unreached []interface{}, quit chan bool) *Client {
	return &Client{conn: conn, Token: token, Unreached: unreached, quit: quit}
}

func (c *Client) Conn() *websocket.Conn {
	c.Lock()
	defer c.Unlock()
	return c.conn
}

func (c *Client) UpdateConn(conn *websocket.Conn){
	c.Lock()
	defer c.Unlock()
	c.conn = conn
}

func (c *Client) SetRoom(r *werewolves.Room){
	c.Lock()
	defer c.Unlock()
	c.room = r
}
func (c *Client) SetChar(plr *werewolves.Player){
	c.Lock()
	defer c.Unlock()
	c.Player = plr
}

func (c *Client) Room()*werewolves.Room{
	c.Lock()
	defer c.Unlock()
	return c.room
}
func (c *Client) Char() *werewolves.Player{
	c.Lock()
	defer c.Unlock()
	return c.Player
}

func (c *Client) ListenRoom(){
	for {
		if c.Player != nil {
			select {
			case value, ok := <-c.Player.Update:
				if ok {
					c.WriteJSON(value)
				} else {
					return
				}
			case <- c.quit:
				return
			}
		}else{
			return
		}
	}
}

func (c *Client) WriteJSON(i interface{}) error {
	c.Lock()
	defer c.Unlock()
	var err error
	if c.conn != nil {
		err = c.conn.WriteJSON(i)
	}
	if err != nil{
		c.Unreached = append(c.Unreached, i)
		log.Println(c.Token.Username, "unreached", i)
	}
	return err
}

func (c *Client) SendUnreached(){
	if len(c.Unreached) > 0{
		for _, e := range c.Unreached{
			c.WriteJSON(e)
			log.Println("sent unreached to", c.Token.Username, "|", e)
		}
	}
	c.Unreached = make([]interface{}, 0)
}
