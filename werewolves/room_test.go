package werewolves

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"testing"
)

func TestRoom_defineRoles(t *testing.T) {
	players := Players{}
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		p := NewPlayer(s, "1")
		go func() {
			select {
			case <-p.Update:
				b, err := json.Marshal(p)
				if err != nil {
					log.Println(err)
				}
				log.Println(string(b))
			}
		}()
		players[s] = p
	}

	room := NewRoom("1", "2", players, Settings{}, "1")
	err := room.Run()
	if err != nil {
		panic(err)
	}
	for _, v := range players {
		log.Printf("PlayerID: %s Role: %v", v.ID, reflect.TypeOf(v.Character))
	}
}