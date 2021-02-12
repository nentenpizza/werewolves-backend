package storage

import (
	"time"

	//sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(User) error
		Exists(username string) (has bool ,_ error )
	}

	Users struct {
		*sqlx.DB
	}

	// User represents user-account
	User struct {
		CreatedAt         time.Time `sq:"created_at,omitempty" json:"created_at"`
		UpdatedAt time.Time `sq:"updated_at,omitempty" json:"updated_at""`
		XP int64 `db:"xp" sq:"xp,omitempty" json:"xp"`
		ID                int    `db:"id" sq:"id" json:"id,omitempty" validate:"required,id"`
		Email             string    `sq:"email" json:"email,omitempty" validate:"required,email"`
		Username          string    `sq:"username" json:"username" validate:"required"`
		EncryptedPassword string    `db:"password_hash" sq:"password" json:"-"`
		Wins int `sq:"wins" json:"wins"`
		Losses int `sq:"losses" json:"losses"`
	}
)

func (db *Users) Create(u User) error {
	const q = `insert into users (email, username, password_hash) values ($1,$2,$3)`
	_, err := db.Exec(q, u.Email, u.Username, u.EncryptedPassword)
	return err
}

func (db *Users) Exists(username string) (has bool ,_ error ) {
	const q = `select exists(select * from users where username = $1)`
	return has, db.Get(&has, q, username)
}