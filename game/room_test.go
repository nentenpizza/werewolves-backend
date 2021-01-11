package game

import (
	"log"
	"reflect"
	"strconv"
	"testing"
)

func TestRoom_defineRoles(t *testing.T) {
	players := Players{}
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		players[s] = NewPlayer(s)
	}

	room := NewRoom(players)
	room.Run()
	for _, v := range players {
		log.Println(reflect.TypeOf(v.Role))
	}
}
