package storage

import (
	"github.com/lib/pq"
	"time"

	//sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(User) error
		Exists(username string) (has bool, _ error)
		ByUsername(username string) (user User, _ error)
		ByID(id int64) (user User, _ error)
		ExistsByID(id int64) (exists bool, _ error)
		ExistsByLogin(login string) (has bool, _ error)
		Relations(User) (pq.Int64Array, error)
		ByLogin(login string) (user User, _ error)
		Update(user User) error
		UpdateRelations(userID int64, relID int) error
	}

	Users struct {
		*sqlx.DB
	}

	// User represents user-account
	User struct {
		CreatedAt         time.Time     `sq:"created_at,omitempty" json:"created_at"`
		UpdatedAt         time.Time     `sq:"updated_at,omitempty" json:"-"`
		XP                int64         `db:"xp" sq:"xp,omitempty" json:"xp"`
		ID                int64         `db:"id" sq:"id" json:"id,omitempty" validate:"required,id"`
		Email             string        `sq:"email" json:"email,omitempty" validate:"required,email"`
		Login             string        `sq:"login" json:"-" validate:"required"`
		Username          string        `sq:"username" json:"username" validate:"required"`
		Relations         pq.Int64Array `sq:"relations" json:"relations"`
		EncryptedPassword string        `db:"password_hash" sq:"password" json:"-"`
		BannedUntil       time.Time     `sq:"banned_until,omitempty" json:"banned_until"`
		Avatar            string        `sq:"avatar" json:"avatar"`
		Wins              int           `sq:"wins" json:"wins"`
		Losses            int           `sq:"losses" json:"losses"`
	}
)

// Sanitize removes User.Email, User.EncryptedPassword, User.BannedUntil from struct
func (u *User) Sanitize() {
	u.Email = ""
	u.EncryptedPassword = ""
	u.BannedUntil = time.Time{}
}

func (db *Users) Create(u User) error {
	const q = `insert into users (email, username, password_hash, login) values ($1,$2,$3, $4)`
	_, err := db.Exec(q, u.Email, u.Username, u.EncryptedPassword, u.Login)
	return err
}

func (db *Users) Exists(username string) (has bool, _ error) {
	const q = `select exists(select * from users where username = $1)`
	return has, db.Get(&has, q, username)
}

func (db Users) ByUsername(username string) (user User, _ error) {
	const q = `select * from users where username = $1`
	return user, db.Get(&user, q, username)
}

func (db *Users) ExistsByLogin(login string) (has bool, _ error) {
	const q = `select exists(select * from users where login = $1)`
	return has, db.Get(&has, q, login)
}

func (db Users) ByLogin(login string) (user User, _ error) {
	const q = `select * from users where login = $1`
	return user, db.Get(&user, q, login)
}

func (db Users) ExistsByID(id int64) (exists bool, _ error) {
	const q = `select exists(select * from users where id = $1)`
	return exists, db.Get(&exists, q, id)
}

func (db Users) ByID(id int64) (user User, _ error) {
	const q = `select * from users where id = $1`
	return user, db.Get(&user, q, id)
}

func (db Users) Update(user User) error {
	const q = `
		UPDATE users SET 
			xp = :xp,
			wins = :wins,
			losses = :losses,
			avatar = :avatar
		WHERE id = :id`
	_, err := db.NamedExec(q, user)
	return err
}

func (db Users) Relations(user User) (r pq.Int64Array, _ error) {
	const q = "SELECT relations FROM users WHERE id = $1"
	return r, db.Get(&r, q, user.ID)
}

func (db Users) UpdateRelations(userID int64, relID int) error {
	const q = "UPDATE users SET relations = array_append(relations, $1) WHERE id = $2"
	_, err := db.Exec(q, relID, userID)
	return err
}
