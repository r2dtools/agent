package router

import "fmt"

// Handler is an interface that must be implemented by any action handler
type Handler interface {
	Handle(request Request) (interface{}, error)
}

// Router contains handlers to handle requests
type Router struct {
	handlers map[string]Handler
}

// RegisterHandler register module handler
func (r *Router) RegisterHandler(module string, handler Handler) {
	if r.handlers == nil {
		r.handlers = make(map[string]Handler)
	}

	r.handlers[module] = handler
}

// GetHandler returns module handler
func (r *Router) GetHandler(request Request) Handler {
	if r.handlers == nil {
		r.handlers = make(map[string]Handler)
	}

	return r.handlers[request.GetModule()]
}

// HandleRequest handles request
func (r *Router) HandleRequest(request Request) (interface{}, error) {
	var handler Handler

	handler = r.GetHandler(request)

	if handler == nil {
		return nil, fmt.Errorf("could not find handler for command '%s'", request.Command)
	}

	return handler.Handle(request)
}
