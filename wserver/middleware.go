package wserver

type MiddlewareFunc func(HandlerFunc) HandlerFunc
