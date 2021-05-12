package service

import (
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type HonorService interface {
	Honor(fromID int64, toID int64, reason string) error
}

type Honors struct {
	Service
}

func (s Honors) Honor(fromID int64, toID int64, reason string) error {
	exists, err := s.db.Users.ExistsByID(toID)
	if err != nil {
		return err
	}
	if !exists {
		return serviceError(http.StatusBadRequest, "user does not exist")
	}

	if toID == fromID {
		return serviceError(http.StatusBadRequest, "you cannot report yourself")
	}

	err = s.db.Honors.Create(storage.Honor{
		HonoredID: toID,
		Reason:    reason,
		SenderID:  fromID,
	})
	return err
}
