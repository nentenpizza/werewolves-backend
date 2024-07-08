package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	FriendsStorage interface {
		Create(userID int64) (id int, _ error)
		UsersByIDs(relationIDs pq.Int64Array, userID int64) (users []User, _ error)
		Accept(userID int64, relID int) error
		IsFriend(relationIDs pq.Int64Array, friendID int64) (has bool, _ error)
		UnacceptedUsersByIDs(relationIDs pq.Int64Array, userID int64) (users []User, _ error)
		AcceptBySenderID(relationIDs pq.Int64Array, senderID int64, userID int64) error
	}

	Friends struct {
		*sqlx.DB
	}

	UnacceptedUsers struct {
		RequestID int  `json:"user" `
		User      User `json:"user"`
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

// UsersByIDs returns all users under specific relation ids in relationships table
// ignores user with specific user_id
func (db *Friends) UsersByIDs(relationIDs pq.Int64Array, userID int64) (users []User, _ error) {
	const q = `
		select * from users where id in 
		(select user_id from relationship where $2 in (select user_id from relationship where id = any (?)))
		and users.id != $2
	`
	query, args, err := sqlx.In(q, relationIDs, userID)
	if err != nil {
		return nil, err
	}
	return users, db.Select(&users, db.Rebind(query), args...)
}

func (db *Friends) UnacceptedUsersByIDs(relationIDs pq.Int64Array, userID int64) (users []User, _ error) {
	const q = `
		select * from users where id in 
		(select user_id from relationship where $2 not in (select user_id from relationship where id = any (?)))
		and users.id != $2
	`
	query, args, err := sqlx.In(q, relationIDs, userID)
	if err != nil {
		return nil, err
	}
	return users, db.Select(&users, db.Rebind(query), args...)
}

func (db *Friends) IsFriend(relationIDs pq.Int64Array, friendID int64) (has bool, _ error) {
	const q = `select 
       exists(select user_id from relationship where $2 in (select user_id from relationship where id = any(?)));
	`
	query, args, err := sqlx.In(q, relationIDs, friendID)
	if err != nil {
		return false, err
	}
	return has, db.Get(&has, db.Rebind(query), args...)
}

func (db *Friends) AcceptBySenderID(relationIDs pq.Int64Array, senderID int64, userID int64) error {
	const q = `insert into relationship (user_id, id) values ($3, (select id from relationship where $2 in (select user_id from relationship where id = any(?)) );
	`
	query, args, err := sqlx.In(q, relationIDs, senderID, userID)
	if err != nil {
		return err
	}
	_, err = db.Exec(db.Rebind(query), args...)
	return err
}
