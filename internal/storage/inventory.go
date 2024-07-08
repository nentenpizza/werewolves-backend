package storage

import "github.com/jmoiron/sqlx"

type (
	InventoriesStorage interface {
		Create(inv Inventory) error
		Delete(inv Inventory, limit uint) error
		Fetch(UserID int) ([]Inventory, error)
		CountItem(Inventory) (int, error)
	}

	Inventories struct {
		*sqlx.DB
	}

	Inventory struct {
		ID     int    `json:"-" sq:"id"`
		UserID int    `json:"user_id" sq:"user_id"`
		Item   string `json:"item" sq:"item"`
	}
)

func (db Inventories) Create(inv Inventory) error {
	const q = "INSERT INTO inventory (user_id, item) VALUES ($1, $2)"
	_, err := db.Exec(q, inv.UserID, inv.Item)
	return err
}

func (db Inventories) Fetch(UserID int) (i []Inventory, _ error) {
	const q = "SELECT * FROM inventory WHERE user_id = $1"
	return i, db.Select(&i, q, UserID)
}

func (db Inventories) CountItem(inv Inventory) (count int, _ error) {
	const q = "SELECT COUNT(*) FROM inventory WHERE item = $1 AND user_id = $2"
	return count, db.Get(&count, q, inv.Item, inv.UserID)
}

func (db Inventories) Delete(inv Inventory, limit uint) error {
	const q = "DELETE FROM inventory WHERE ctid IN (SELECT ctid FROM inventory WHERE item = $1 AND user_id = $2 LIMIT $3)"
	_, err := db.Exec(q, inv.Item, inv.UserID, limit)
	return err
}
