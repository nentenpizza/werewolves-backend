// app for testing server
package main

import (
	"github.com/nentenpizza/werewolves/handler"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/nentenpizza/werewolves/werewolves"
)

func main() {
	dialer := &websocket.Dialer{}
	ws, _, err := dialer.Dial("ws://localhost:7070/ws", http.Header{})
	if err != nil {
		panic(err)
	}
	ev := &handler.Event{
		handler.EventTypeCreateRoom,
		&handler.EventCreateRoom{"debil", "dura", werewolves.Settings{}},
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
