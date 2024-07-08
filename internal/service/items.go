package service

import (
	"github.com/nentenpizza/werewolves/internal/storage"
)

type ItemService interface {
	Count(userID int64, itemName string) (int, error)
	DeleteItem(userID int64, itemName string, count uint) (int, error)
	GiveItem(userID int64, itemName string) (int, error)
}

type Items struct {
	Service
}

func (s Items) GiveItem(userID int64, itemName string) (int, error) {
	user, err := s.db.Users.ByID(userID)
	if err != nil {
		return 0, err
	}

	item := storage.Item{
		Name:   itemName,
		UserID: user.ID,
	}

	err = s.db.Items.Create(item)
	if err != nil {
		return 0, err
	}
	return s.db.Items.Count(item)
}

func (s Items) Count(userID int64, itemName string) (int, error) {
	user, err := s.db.Users.ByID(userID)

	if err != nil {
		return 0, err
	}

	return s.db.Items.Count(storage.Item{
		Name:   itemName,
		UserID: user.ID,
	})
}

func (s Items) DeleteItem(userID int64, itemName string, count uint) (int, error) {
	user, err := s.db.Users.ByID(userID)

	if err != nil {
		return 0, err
	}

	item := storage.Item{
		UserID: user.ID,
		Name:   itemName,
	}

	err = s.db.Items.Delete(item, count)
	if err != nil {
		return 0, err
	}

	return s.db.Items.Count(item)
}
