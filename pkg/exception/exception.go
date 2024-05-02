package exception

import (
	"fmt"
)

// Exceptions.
var (
	ErrUnauthorized        error = fmt.Errorf("Unauthorized")
	ErrNotFound            error = fmt.Errorf("Not found")
	ErrInternalServer      error = fmt.Errorf("Internal server error")
	ErrConflict            error = fmt.Errorf("Conflict")
	ErrUnprocessableEntity error = fmt.Errorf("Unprocessable entity")
	ErrBadRequest          error = fmt.Errorf("Bad request")
	ErrGatewayTimeout      error = fmt.Errorf("Gateway timeout")
	ErrTimeout             error = fmt.Errorf("Request time out")
	ErrLocked              error = fmt.Errorf("Locked")
	ErrForbidden           error = fmt.Errorf("Forbidden")
)
