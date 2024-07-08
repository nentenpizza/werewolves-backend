package service

import "net/http"

var (
	InvalidUsername = serviceError(http.StatusBadRequest, "usr 3-10 chars")
	InvalidLogin    = serviceError(http.StatusBadRequest, "login 3-16 chars")
	InvalidPassword = serviceError(http.StatusBadRequest, "password 5-30 chars")
)
