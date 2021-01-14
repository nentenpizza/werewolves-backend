// app for testing server
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/nentenpizza/werewolves/werewolves"

	"github.com/nentenpizza/werewolves/server"
)

func main() {
	dialer := &websocket.Dialer{}
	ws, _, err := dialer.Dial("ws://localhost:7070/ws", http.Header{})
	if err != nil {
		panic(err)
	}
	ev := &server.Event{
		server.EventTypeCreateRoom,
		&server.EventCreateRoom{"debil", "dura", werewolves.Settings{}},
	}
	if err != nil {
		panic(err)
	}

	for {
		if err := ws.WriteJSON(ev); err != nil {
			panic(err)
		}
		_, msg, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		log.Println(string(msg))
	}
}
