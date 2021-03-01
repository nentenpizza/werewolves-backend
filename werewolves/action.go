package werewolves

var allowedActions = map[string][]string{
	Discuss:   {ConstableShootAction, DoctorHealAction},
	DayVoting: {ConstableShootAction, DoctorHealAction, VoteAction},
	Night:     {DoctorHealAction},
}

// Action names
const (
	ConstableShootAction   = "constable.shoot"
	DoctorHealAction       = "doctor.heal"
	VoteAction             = "game.vote"
	PsychicResurrectAction = "psychic.resurrect"
	ExecutionAction        = "game.execution"
)

// Action represents players actions
// for example shoot or heal
// Room performs actions
type Action struct {
	Name string

	// (do) must return an error if action cannot be performed
	do func(r *Room) error

	Event Event
}

// NewAction creates Action
func NewAction(name string, do func(r *Room) error, ev Event) Action {
	return Action{Name: name, do: do, Event: ev}
}
