package wserver

// MiddlewareFunc represents a middleware processing function,
// which get called before the endpoint group or specific handler.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Group is a separated group of handlers, united by the general middleware.
type Group struct {
	s          *Server
	middleware []MiddlewareFunc
}

// Use adds middleware to the chain.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Handle adds endpoint handler to the server, combining group's middleware
// with the optional given middleware.
func (g *Group) Handle(endpoint interface{}, h HandlerFunc, m ...MiddlewareFunc) {
	g.s.Handle(endpoint, h, append(g.middleware, m...)...)
}
