package response

// Response is a collection of response's behavior.
type Response interface {
	Data() interface{}
	Error() error
	Status() string
	HTTPStatusCode() int
	Message() string
	Meta() interface{}
}
