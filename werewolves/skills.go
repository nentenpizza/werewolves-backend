package werewolves

import (
	"errors"
	"log"
)

func (char *Constable) Shoot(other *Player) Action {
	char.Lock()
	defer char.Unlock()
	return NewAction(
		ConstableShootAction,

		func(_ *Room) error {
			if char.HP() <= 0{
				return errors.New("werewolves: player are dead")
			}
			if other.Character.HP() <= 0{
				return errors.New("werewolves: target already dead")
			}
			if char.bullets <= 0 {
				return errors.New("game: constable out of bullets")
			}
			other.Character.SetHP(other.Character.HP() - 1)
			log.Println("Shoot in ", other.ID)
			return nil
		},

		NewEvent(ConstableShootAction, ConstableShootEvent{PlayerID: char.ParentID, TargetID: other.ID}),
	)
}

func (char *Doctor) Heal(other *Player) Action {
	char.Lock()
	defer char.Unlock()
	return NewAction(DoctorHealAction,
		func(_ *Room) error {
			if char.HP() <= 0{
				return errors.New("werewolves: player are dead")
			}
			other.Character.SetHP(other.Character.HP() + 1)
			log.Println("Heal player ", other.ID)
			return nil
		},
		NewEvent(
			PsychicResurrectAction,
			TargetedEvent{other.ID},
		),
	)
}

func (char *Psychic) Resurrect(other *Player) Action {
	char.Lock()
	defer char.Unlock()
	return NewAction(
		PsychicResurrectAction,
		func(r *Room) error {
			if char.HP() <= 0{
				return errors.New("werewolves: player are dead")
			}
			err := r.Resurrect(other.ID)
			return err
		},
		NewEvent(
			PsychicResurrectAction,
			TargetedEvent{other.ID},
		),
	)

}
