package cmd

import (
	"encoding/json"
	"github.com/nentenpizza/werewolves-backend/game"
	"log"
	"reflect"
	"strconv"
)

func main() {
	players := game.Players{}
	for i := 0; i < 2; i++ {
		s := strconv.Itoa(i)
		p := game.NewPlayer(s)
		go func() {
			select {
			case <-p.Update:
				b, err := json.Marshal(game.NewPlayerState(p))
				if err != nil {
					log.Println(err)
				}
				log.Println(string(b))
			}
		}()
		players[s] = p
	}

	room := game.NewRoom(players, Settings{})
	err := room.Run()
	if err != nil {
		panic(err)
	}
	constable := players["0"].Character.(*Constable)
	err = room.Perform(constable.Shoot(players["1"]))
	if err != nil {
		panic(err)
	}
	for _, v := range players {
		log.Printf("PlayerID: %s Role: %v", v.ID, reflect.TypeOf(v.Character))
	}
}
