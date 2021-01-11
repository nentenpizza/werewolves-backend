package game

var rolesMap = map[int][]func() Character{
	1:  []func() Character{newConstable},
	2:  []func() Character{newConstable, newConstable},
	3:  []func() Character{newConstable, newConstable, newConstable},
	4:  []func() Character{newConstable, newConstable, newConstable, newConstable},
	5:  []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable},
	6:  []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable, newConstable},
	7:  []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable},
	8:  []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable},
	9:  []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable},
	10: []func() Character{newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable, newConstable},
}

// Character represents interface for each role in game
// All roles in game must satisfy this interface.
// By default all roles must have 1hp but some roles can have more
type Character interface {
	HP() int
	SetHP(hp int)
	Died() bool
}

type Constable struct {
	hp int
}

func (char *Constable) Died() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newConstable() Character {
	return &Constable{hp: 1}
}

func (char *Constable) SetHP(hp int) {
	char.hp = char.hp + hp

}

func (char *Constable) HP() int {
	return char.hp
}
