package middleware

import "net/http"

// RouteMiddleware is a abstraction of route middleware.
type RouteMiddleware interface {
	Verify(next http.HandlerFunc) http.HandlerFunc
}

// RecaptchaRouteMiddleware is a abstraction of recaptcha route middleware.
type RecaptchaRouteMiddleware interface {
	Verify(next http.HandlerFunc, actions ...string) http.HandlerFunc
}
