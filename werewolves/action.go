package werewolves

var allowedActions = map[string][]string{
	Discuss:   []string{ConstableShootAction, DoctorHealAction},
	DayVoting: []string{ConstableShootAction, DoctorHealAction},
	Night:     []string{DoctorHealAction},
}

// Action names
const (
	ConstableShootAction = "constable.shoot"
	DoctorHealAction     = "doctor.heal"
	VoteAction           = "game.vote"
)

// Action represents players actions
// for example shoot or heal
// Room performs actions
type Action struct {
	Name string
	do   func(r *Room) error
}

// NewAction creates Action
func NewAction(name string, do func(r *Room) error) Action {
	return Action{Name: name, do: do}
}
