package server

import "net/http"

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
}
