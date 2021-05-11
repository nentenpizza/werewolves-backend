package service

import (
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type ReportService interface {
	Report(fromID int64, toID int64, reason string) error
}

type Reports struct {
	Service
}

func (s Reports) Report(fromID int64, toID int64, reason string) error {
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

	err = s.db.Reports.Create(storage.Report{
		ReportedID: toID,
		Reason:     reason,
		SenderID:   fromID,
	})
	return err
}
