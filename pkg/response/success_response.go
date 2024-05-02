package response

import "net/http"

// SuccessResponse is an model of success response.
type SuccessResponse struct {
	httpStatusCode int
	status         string
	message        string
	data           interface{}
	meta           interface{}
}

// NewSuccessResponse is a constructor.
func NewSuccessResponse(data interface{}, status string, message string) Response {
	var httpStatusCode int
	switch status {
	case StatCreated:
		httpStatusCode = http.StatusCreated
		break
	default:
		httpStatusCode = http.StatusOK
		break
	}

	return SuccessResponse{
		httpStatusCode: httpStatusCode,
		status:         status,
		message:        message,
		data:           data,
	}
}

// NewSuccessResponseWithMeta is a constructor.
func NewSuccessResponseWithMeta(data interface{}, meta interface{}, status string, message string) Response {
	var httpStatusCode int
	switch status {
	case StatCreated:
		httpStatusCode = http.StatusCreated
		break
	default:
		httpStatusCode = http.StatusOK
		break
	}

	return SuccessResponse{
		httpStatusCode: httpStatusCode,
		status:         status,
		message:        message,
		data:           data,
		meta:           meta,
	}
}

// Data returns data.
func (r SuccessResponse) Data() interface{} {
	return r.data
}

// Error return error.
func (r SuccessResponse) Error() error {
	return nil
}

// Status returns response status.
func (r SuccessResponse) Status() string {
	return r.status
}

// HTTPStatusCode returns http status code.
func (r SuccessResponse) HTTPStatusCode() int {
	return r.httpStatusCode
}

// Message returns message.
func (r SuccessResponse) Message() string {
	return r.message
}

// Meta returns meta.
func (r SuccessResponse) Meta() interface{} {
	return r.meta
}
