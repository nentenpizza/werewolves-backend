package game

var allowedActions = map[string][]string{
	Discuss:   []string{ConstableShootAction},
	DayVoting: []string{ConstableShootAction},
	Night:     []string{},
}

type Action struct {
	Name string
	do   func(r *Room)
}

func NewAction(name string, do func(r *Room)) Action {
	return Action{Name: name, do: do}
}
