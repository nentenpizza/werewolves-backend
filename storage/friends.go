package storage

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	FriendsStorage interface {
		Add(Friend) error
		Delete(Friend) error

		Accept(Friend) error
		Decline(Friend) error

		ByUserID(UserID int, active int) ([]Friend, error)

		IsFriend(TargetID, UserID int) (bool, error)
	}

	Friends struct {
		*sqlx.DB
	}

	// Friend represents friend
	Friend struct {
		ID        int       `json:"-" sq:"id"`
		CreatedAt time.Time `sq:"created_at" json:"created_at"`
		SenderID  int       `json:"sender_id" sq:"sender_id"`
		TargetID  int       `json:"target_id" sq:"target_id"`
		Active    bool      `json:"active" sq:"active"`
	}
)

func (db Friends) Add(f Friend) error {
	const q = "INSERT INTO friends (sender_id, target_id, active) VALUES ($1, $2, $3)"
	_, err := db.Exec(q, f.SenderID, f.TargetID, f.SenderID)
	return err
}

func (db Friends) ByUserID(UserID, active int) (fs []Friend, _ error) {
	const q = "SELECT * FROM friends WHERE sender_id = $1 AND active = $2"
	return fs, db.Select(&fs, q, UserID, active)
}

func (db Friends) Delete(f Friend) error {
	const q = "DELETE FROM friends WHERE sender_id = $1 AND target_id = $2"
	_, err := db.Exec(q, f.SenderID, f.TargetID)
	return err
}

func (db Friends) IsFriend(TargetID, UserID int) (exists bool, _ error) {
	const q = "SELECT EXISTS (SELECT * FROM friends WHERE target_id = $1 AND sender_id = $2 AND active = 0)"
	return exists, db.Get(&exists, q, TargetID, UserID)
}

func (db Friends) Accept(f Friend) error {
	const q = "UPDATE friends SET active = 0 WHERE target_id = $1 AND sender_id = $2"
	_, err := db.Exec(q, f.TargetID, f.SenderID)
	return err
}

func (db Friends) Decline(f Friend) error {
	const q = "DELETE FROM friends WHERE target_id = $1 AND sender_id = $2 AND active = 0"
	_, err := db.Exec(q, f.TargetID, f.SenderID)
	return err
}
