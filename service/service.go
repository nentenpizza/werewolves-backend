package service

import "github.com/nentenpizza/werewolves/storage"

type Service struct {
	db *storage.DB
}

func NewService(db *storage.DB) Service {
	return Service{
		db: db,
	}
}

type Error struct {
	Code    int
	Message string
}

func (s Error) Error() string {
	return s.Message
}

func serviceError(code int, msg string) *Error {
	return &Error{
		Code: code, Message: msg,
	}
}
