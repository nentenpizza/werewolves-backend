package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	FriendsStorage interface {
		Create(userID int64) (id int, _ error)
		UsersByID(relationIDs pq.Int64Array, userID int64) (users []User, _ error)
		Accept(userID int64, relID int) error
	}

	Friends struct {
		*sqlx.DB
	}
)

func (db *Friends) Create(userID int64) (id int, _ error) {
	const q = `insert into relationship (user_id) values ($1) returning id`

	rows, err := db.Query(q, userID)
	if err != nil {
		return 0, err
	}

	rows.Next()

	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return
}

func (db *Friends) Accept(userID int64, relID int) error {
	const q = `insert into relationship (user_id, id) values ($1, $2)`
	_, err := db.Exec(q, userID, relID)
	return err
}

// UsersByID returns all users under specific relation ids in relationships table
// ignores user with specific user_id
func (db *Friends) UsersByID(relationIDs pq.Int64Array, userID int64) (users []User, _ error) {
	const q = `
		select * from users where id in 
		(select user_id from relationship where id = any (?))
		and users.id != $2
	`
	query, args, err := sqlx.In(q, relationIDs, userID)
	if err != nil {
		return nil, err
	}
	return users, db.Select(&users, db.Rebind(query), args...)
}
