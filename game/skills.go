package game

const (
	ConstableShootAction = "constable.shoot"
)

func (char *Constable) Shoot(other Character) Action {
	return NewAction(ConstableShootAction, func(_ *Room) {
		other.SetHP(-1)
	})
}
