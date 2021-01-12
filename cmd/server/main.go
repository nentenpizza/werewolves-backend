package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nentenpizza/werewolves/game"

	"github.com/gorilla/websocket"
)

func main() {

}

type Event struct {
	Type string `json:"type"`
}

// serveWs handles websocket requests from the peer.
func serveWs(room *game.Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			b, err := json.Unmarshal()
		}
	}()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
