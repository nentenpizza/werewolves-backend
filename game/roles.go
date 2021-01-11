package game

var rolesMap = map[int][]func() Character{
	1:  {newConstable},
	2:  {newConstable, newConstable},
	3:  {newVillager, newConstable, newWerewolf},
	4:  {newConstable, newWerewolf, newVillager, newVillager},
	5:  {newVillager, newWerewolf, newVillager, newVillager, newConstable},
	6:  {newVillager, newVillager, newConstable, newWerewolf, newVillager, newVillager},
	7:  {newConstable, newConstable, newWerewolf, newVillager, newVillager, newVillager, newConstable},
	8:  {newWerewolf, newConstable, newWerewolf, newVillager, newVillager, newVillager, newVillager, newConstable},
	9:  {newWerewolf, newDoctor, newWerewolf, newPsychic, newVillager, newVillager, newVillager, newVillager, newConstable},
	10: {newAlphaWerewolf, newWerewolf, newWerewolf, newPsychic, newFool, newVillager, newVillager, newVillager, newVillager, newConstable},
}

// Character represents interface for each role in game
// All roles in game must satisfy this interface.
// By default all roles must have 1hp but some roles can have more
type Character interface {
	HP() int
	SetHP(hp int)
	Dead() bool
}

// Constable role
type Constable struct {
	hp int
}

func (char *Constable) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newConstable() Character {
	return &Constable{hp: 1}
}

func (char *Constable) SetHP(hp int) {
	char.hp = hp

}

func (char *Constable) HP() int {
	return char.hp
}

// Werewolf role
type Werewolf struct {
	hp int
}

func (char *Werewolf) HP() int {
	return char.hp
}

func (char *Werewolf) SetHP(hp int) {
	char.hp = hp
}

func newWerewolf() Character {
	return &Werewolf{hp: 1}
}

func (char *Werewolf) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

// AlphaWerewolf role
type AlphaWerewolf struct {
	hp int
}

func (char *AlphaWerewolf) HP() int {
	return char.hp
}

func (char *AlphaWerewolf) SetHP(hp int) {
	char.hp = hp
}

func (char *AlphaWerewolf) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newAlphaWerewolf() Character {
	return &AlphaWerewolf{hp: 1}
}

// Doctor role
type Doctor struct {
	hp int
}

func (char *Doctor) HP() int {
	return char.hp
}

func (char *Doctor) SetHP(hp int) {
	char.hp = hp
}

func (char *Doctor) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newDoctor() Character {
	return &Doctor{hp: 1}
}

// Psychic role
type Psychic struct {
	hp int
}

func (char *Psychic) HP() int {
	return char.hp
}

func (char *Psychic) SetHP(hp int) {
	char.hp = hp
}

func (char *Psychic) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newPsychic() Character {
	return &Psychic{hp: 1}
}

// Villager role
type Villager struct {
	hp int
}

func (char *Villager) HP() int {
	return char.hp
}

func (char *Villager) SetHP(hp int) {
	char.hp = hp
}

func (char *Villager) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newVillager() Character {
	return &Villager{hp: 1}
}

// Fool role
type Fool struct {
	hp int
}

func (char *Fool) HP() int {
	return char.hp
}

func (char *Fool) SetHP(hp int) {
	char.hp = hp
}

func (char *Fool) Dead() bool {
	if char.hp <= 0 {
		return true
	}
	return false
}

func newFool() Character {
	return &Fool{hp: 1}
}
