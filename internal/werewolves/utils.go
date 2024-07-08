package werewolves

import (
	werewolves "github.com/nentenpizza/werewolves/pkg/werewolves"
)

func (g game) Vote(owner *werewolves.Player, target string) error {
	var vCount uint8 = 1
	if owner.Role == "AlphaWerewolf" && owner.Room.State == werewolves.Night {
		vCount = 2
	}
	action := owner.Vote(target, vCount)
	if owner.Room.State == werewolves.Night {
		if err := owner.Room.Perform(action, "wolves"); err != nil {
			return NotAllowedErr
		}
	} else {
		if err := owner.Room.Perform(action); err != nil {
			return NotAllowedErr
		}
	}
	return nil
}
