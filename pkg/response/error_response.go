package response

// ErrorResponse is an model of success response.
type ErrorResponse struct {
	err            error
	httpStatusCode int
	status         string
	message        string
	data           interface{}
}

// NewErrorResponse is a constructor.
func NewErrorResponse(err error, httpStatusCode int, data interface{}, status string, message string) Response {
	return ErrorResponse{
		err:            err,
		httpStatusCode: httpStatusCode,
		status:         status,
		message:        message,
		data:           data,
	}
}

// Data returns data.
func (r ErrorResponse) Data() interface{} {
	return r.data
}

// Error returns error.
func (r ErrorResponse) Error() error {
	return r.err
}

// Status returns response status.
func (r ErrorResponse) Status() string {
	return r.status
}

// HTTPStatusCode returns http status code.
func (r ErrorResponse) HTTPStatusCode() int {
	return r.httpStatusCode
}

// Message returns message.
func (r ErrorResponse) Message() string {
	return r.message
}

// Meta reutrns meta.
func (r ErrorResponse) Meta() interface{} {
	return nil
}
