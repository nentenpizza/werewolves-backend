package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	HonorsStorage interface {
		Create(Honor) error
		Exists(honored, sender int64) (bool, error)
		CountToday(UserID int64) (c int, _ error)
	}

	Honors struct {
		*sqlx.DB
	}

	// Honor represents honor
	Honor struct {
		ID        int       `json:"id" sq:"id"`
		CreatedAt time.Time `json:"created_at" sq:"created_at"`
		HonoredID int64     `json:"honored_id" sq:"honored_id"`
		SenderID  int64     `json:"sender_id" sq:"sender_id"`
		Reason    string    `json:"reason" sq:"reason"`
	}
)

func (db Honors) Create(h Honor) error {
	const q = "INSERT INTO honors (honored_id, sender_id, reason) VALUES ($1, $2, $3)"
	_, err := db.Exec(q, h.HonoredID, h.SenderID, h.Reason)
	return err
}

func (db Honors) Exists(HonoredID, SenderID int64) (e bool, _ error) {
	const q = "SELECT EXISTS(SELECT * FROM honors WHERE honored_id = $1 AND sender_id = $2)"
	return e, db.Get(&e, q, HonoredID, SenderID)
}

func (db Honors) CountToday(UserID int64) (c int, _ error) {
	const q = "SELECT COUNT(*) FROM honors WHERE sender_id = $1 AND created_at::date between date 'now()' and date 'now()'"
	return c, db.Get(&c, q, UserID)
}
