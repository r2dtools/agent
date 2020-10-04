package server

import "strings"

// Request is a request object from the main server
type Request struct {
	Command,
	Token string
	Data interface{}
}

// GetModule returns module that should hanlde the request
func (r *Request) GetModule() string {
	parts := strings.SplitN(r.Command, ".", 2)

	if len(parts) == 2 {
		return parts[0]
	}

	return ""
}

// GetAction returns action to handle
func (r *Request) GetAction() string {
	parts := strings.SplitN(r.Command, ".", 2)

	if len(parts) == 2 {
		return parts[1]
	}

	return r.Command
}
