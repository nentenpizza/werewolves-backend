package game

import (
	"log"
)

const (
	ConstableShootAction = "constable.shoot"
	DoctorHealAction     = "doctor.heal"
)

func (char *Constable) Shoot(other *Player) Action {
	return NewAction(ConstableShootAction, func(_ *Room) {
		other.Character.SetHP(other.Character.HP() - 1)
		other.Room.Broadcast <- true
		log.Println("Shoot in ", other.ID)
	})
}
func (char *Doctor) Heal(other *Player) Action {
	return NewAction(DoctorHealAction, func(_ *Room) {
		other.Character.SetHP(other.Character.HP() + 1)
		other.Room.Broadcast <- true
		log.Println("Heal player ", other.ID)
	})
}
