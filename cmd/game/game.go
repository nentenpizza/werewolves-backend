package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nentenpizza/werewolves/werewolves"
)

//
func main() {
	players := werewolves.Players{}
	room := werewolves.NewRoom("1", "sample_room", players, werewolves.Settings{OpenRolesOnDeath: true})
	for i := 0; i < 2; i++ {
		s := strconv.Itoa(i)
		p := werewolves.NewPlayer(s)
		go func() {
			for {
				select {
				case <-p.Update:
					b, err := json.Marshal(p)
					if err != nil {
						log.Println(err)
					}
					log.Println(string(b) + strings.Repeat("\n----------------------------", 2))
				}
			}
		}()
		players[s] = p
	}

	err := room.Run()

	if err != nil {
		panic(err)
	}

	constable, ok := players["0"].Character.(*werewolves.Constable)
	if ok {
		room.Perform(constable.Shoot(players["1"]))
	} else {

		doctor := players["0"].Character.(*werewolves.Doctor)
		room.Perform(doctor.Heal(players["0"]))
		constable := players["1"].Character.(*werewolves.Constable)
		room.Perform(constable.Shoot(players["0"]))
	}

	for _, v := range players {
		log.Printf("PlayerID: %s Role: %v", v.ID, reflect.TypeOf(v.Character))
	}
	time.Sleep(123 * time.Second)
}
