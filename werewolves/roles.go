package werewolves

import "sync"

var rolesMap = map[int][]func(parentID string) Character{
	1:  {newConstable},
	2:  {newConstable, newAlphaWerewolf},
	3:  {newDoctor, newPsychic, newWerewolf},
	4:  {newConstable, newWerewolf, newVillager, newVillager},
	5:  {newVillager, newWerewolf, newVillager, newVillager, newConstable},
	6:  {newVillager, newAlphaWerewolf, newVillager, newVillager, newVillager, newVillager},
	7:  {newConstable, newConstable, newWerewolf, newVillager, newVillager, newVillager, newConstable},
	8:  {newWerewolf, newConstable, newWerewolf, newVillager, newVillager, newVillager, newVillager, newConstable},
	9:  {newWerewolf, newDoctor, newWerewolf, newPsychic, newVillager, newVillager, newVillager, newVillager, newConstable},
	10: {newWerewolf, newAlphaWerewolf, newWerewolf, newVillager, newVillager, newVillager, newVillager, newVillager, newVillager, newConstable},
	11: {newVillager, newAlphaWerewolf, newWerewolf, newWerewolf, newVillager, newVillager, newVillager, newVillager, newVillager, newVillager, newVillager},
	12: {newVillager, newAlphaWerewolf, newWerewolf, newWerewolf, newWerewolf, newVillager, newVillager, newVillager, newVillager, newVillager, newVillager, newVillager},
}

type CharSettings struct {
	parentID string
}

// Character represents interface for each role in game
// All roles in game must satisfy this interface.
// By default all roles must have 1Hp but some roles can have more
type Character interface {
	HP() int
	SetHP(Hp int)
	IsDead() bool
}

// Constable role
type Constable struct {
	Hp      int
	Dead    bool
	bullets uint8
	sync.Mutex
	ParentID string
}

func (char *Constable) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}

}

func (char *Constable) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func (char *Constable) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func newConstable(parentID string) Character {
	return &Constable{Hp: 1, bullets: 2, ParentID: parentID}
}

// Werewolf role
type Werewolf struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *Werewolf) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *Werewolf) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func newWerewolf(parentID string) Character {
	return &Werewolf{Hp: 1}
}

func (char *Werewolf) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

// AlphaWerewolf role
type AlphaWerewolf struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *AlphaWerewolf) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *AlphaWerewolf) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func (char *AlphaWerewolf) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func newAlphaWerewolf(parentID string) Character {
	return &AlphaWerewolf{Hp: 1}
}

// Doctor role
type Doctor struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *Doctor) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *Doctor) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func (char *Doctor) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func newDoctor(parentID string) Character {
	return &Doctor{Hp: 1}
}

// Psychic role
type Psychic struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *Psychic) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *Psychic) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func (char *Psychic) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func newPsychic(parentID string) Character {
	return &Psychic{Hp: 1}
}

// Villager role
type Villager struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *Villager) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *Villager) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func (char *Villager) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func newVillager(parentID string) Character {
	return &Villager{Hp: 1}
}

// Fool role
type Fool struct {
	Hp   int
	Dead bool
	sync.Mutex
}

func (char *Fool) HP() int {
	char.Lock()
	defer char.Unlock()
	return char.Hp
}

func (char *Fool) SetHP(Hp int) {
	char.Lock()
	defer char.Unlock()
	char.Hp = Hp
	if char.Hp <= 0 {
		char.Dead = true
	}
}

func (char *Fool) IsDead() bool {
	char.Lock()
	defer char.Unlock()
	return char.Dead
}

func newFool(parentID string) Character {
	return &Fool{Hp: 1}
}
