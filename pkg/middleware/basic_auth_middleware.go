package middleware

import (
	"net/http"

	"pii-encrypt-example/pkg/exception"
	"pii-encrypt-example/pkg/response"
)

const (
	errorMessage = "Invalid token"
)

// BasicAuth is a concrete struct of basic auth verifier.
type BasicAuth struct {
	username, password string
}

// NewBasicAuth is a constructor.
func NewBasicAuth(username, password string) RouteMiddleware {
	return &BasicAuth{username, password}
}

func (ba *BasicAuth) respondUnauthorized(w http.ResponseWriter) {
	resp := response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, errorMessage)
	response.JSON(w, resp)
}

// Verify will verify the request to ensure it comes with an authorized basic auth token.
func (ba *BasicAuth) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			ba.respondUnauthorized(w)
			return
		}

		if !(username == ba.username && password == ba.password) {
			ba.respondUnauthorized(w)
			return
		}
		next(w, r)
	})
}
