package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/nentenpizza/werewolves/game"
)

//
func main() {
	players := game.Players{}
	room := game.NewRoom("1", "sample_room", players, game.Settings{})
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		p := game.NewPlayer(s)
		go func() {
			for {
				select {
				case <-p.Update:
					b, err := json.Marshal(p)
					if err != nil {
						log.Println(err)
					}
					log.Println(string(b) + "\n----------------------------\n")
				}
			}
		}()
		players[s] = p
	}

	err := room.Run()
	if err != nil {
		panic(err)
	}

	for _, v := range players {
		log.Printf("PlayerID: %s Role: %v", v.ID, reflect.TypeOf(v.Character))
	}
	time.Sleep(123 * time.Second)
}
