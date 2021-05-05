package storage

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	ReportsService interface {
		Create(Report) error
	}

	Reports struct {
		*sqlx.DB
	}

	//Report represents report
	Report struct {
		ID         int       `json:"id" sq:"id"`
		CreatedAt  time.Time `json:"created_at" sq:"created_at"`
		ReportedID int       `json:"reported_id" sq:"reported_id"`
		SenderID   int       `json:"sender_id" sq:"sender_id"`
		Reason     string    `json:"reason" sq:"reason"`
	}
)

func (db Reports) Create(r Report) error {
	const q = "INSERT INTO reports (reported_id, sender_id, reason) VALUES ($1, $2, $3)"
	_, err := db.Exec(q, r.ReportedID, r.SenderID, r.Reason)
	return err
}
