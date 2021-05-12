package service

import (
	"github.com/nentenpizza/werewolves/storage"
	"net/http"
)

type FriendService interface {
	Request(fromID int64, toID int64) (int, error)
	UserFriends(userID int64) ([]storage.User, error)
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
	users, err := s.db.Friends.UsersByID(user.Relations, user.ID)
	if err != nil {
		return nil, err
	}
	return users, nil
}
