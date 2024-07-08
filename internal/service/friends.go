package service

import (
	"net/http"

	"github.com/nentenpizza/werewolves/internal/storage"
)

type FriendService interface {
	Request(fromID int64, toID int64) (int, error)
	UserFriends(userID int64) ([]storage.User, error)
	AcceptBySenderID(userID int64, senderID int64) error
	UnacceptedUsers(userID int64) ([]storage.User, error)
}

type Friends struct {
	Service
}

func (s Friends) Request(fromID int64, toID int64) (int, error) {
	if toID == fromID {
		return 0, serviceError(http.StatusBadRequest, "you cannot request yourself")
	}
	me, err := s.db.Users.ByID(fromID)
	if err != nil {
		return 0, err
	}

	has, err := s.db.Friends.IsFriend(me.Relations, toID)
	if err != nil {
		return 0, err
	}

	if has {
		return 0, serviceError(http.StatusConflict, "receiver already your friend")
	}
	id, err := s.db.Friends.Create(fromID)
	if err != nil {
		return 0, err
	}

	err = s.db.Users.UpdateRelations(toID, id)
	if err != nil {
		return 0, err
	}

	err = s.db.Users.UpdateRelations(fromID, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Friends) UserFriends(userID int64) ([]storage.User, error) {
	user, err := s.db.Users.ByID(userID)
	if err != nil {
		return nil, err
	}

	users, err := s.db.Friends.UsersByIDs(user.Relations, user.ID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s Friends) AcceptBySenderID(userID int64, senderID int64) error {
	me, err := s.db.Users.ByID(userID)
	if err != nil {
		return err
	}

	has, err := s.db.Friends.IsFriend(me.Relations, senderID)
	if err != nil {
		return err
	}

	if has {
		return serviceError(http.StatusConflict, "user already your friend")
	}
	err = s.db.Friends.AcceptBySenderID(me.Relations, senderID, userID)
	return err
}

func (s Friends) UnacceptedUsers(userID int64) ([]storage.User, error) {
	me, err := s.db.Users.ByID(userID)
	if err != nil {
		return nil, err
	}

	return s.db.Friends.UnacceptedUsersByIDs(me.Relations, userID)
}
