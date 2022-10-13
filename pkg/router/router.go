package router

import "fmt"

type Handler interface {
	Handle(request Request) (interface{}, error)
}

type Router struct {
	handlers map[string]Handler
}

func (r *Router) RegisterHandler(module string, handler Handler) {
	if r.handlers == nil {
		r.handlers = make(map[string]Handler)
	}

	r.handlers[module] = handler
}

func (r *Router) GetHandler(request Request) Handler {
	if r.handlers == nil {
		r.handlers = make(map[string]Handler)
	}

	return r.handlers[request.GetModule()]
}

func (r *Router) HandleRequest(request Request) (interface{}, error) {
	handler := r.GetHandler(request)

	if handler == nil {
		return nil, fmt.Errorf("could not find handler for the command '%s'", request.Command)
	}

	return handler.Handle(request)
}
