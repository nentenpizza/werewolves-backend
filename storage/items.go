package storage

import "github.com/jmoiron/sqlx"

type (
	ItemsStorage interface {
		Create(inv Item) error
		Delete(inv Item, limit uint) error
		Items(UserID int) ([]Item, error)
		Count(Item) (int, error)
	}

	Items struct {
		*sqlx.DB
	}

	Item struct {
		ID     int    `json:"-" sq:"id"`
		UserID int    `json:"user_id" sq:"user_id"`
		Name   string `json:"name" sq:"name"`
	}
)

func (db Items) Create(inv Item) error {
	const q = "INSERT INTO items (user_id, name) VALUES ($1, $2)"
	_, err := db.Exec(q, inv.UserID, inv.Name)
	return err
}

func (db Items) Items(UserID int) (i []Item, _ error) {
	const q = "SELECT * FROM items WHERE user_id = $1"
	return i, db.Select(&i, q, UserID)
}

func (db Items) Count(inv Item) (count int, _ error) {
	const q = "SELECT COUNT(*) FROM items WHERE name = $1 AND user_id = $2"
	return count, db.Get(&count, q, inv.Name, inv.UserID)
}

func (db Items) Delete(inv Item, limit uint) error {
	const q = "DELETE FROM items WHERE ctid IN (SELECT ctid FROM inventory WHERE name = $1 AND user_id = $2 LIMIT $3)"
	_, err := db.Exec(q, inv.Name, inv.UserID, limit)
	return err
}
