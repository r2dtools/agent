package router

// Response that will be sent to the mail server
type Response struct {
	Status,
	Error string
	Data interface{}
}
