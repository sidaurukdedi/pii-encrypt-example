package middleware

import (
	"context"
	"net/http"

	"pii-encrypt-example/entity"
)

func ClientDeviceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		clientDevice := entity.ClientDevice{
			RemoteAddress: r.RemoteAddr,
			XForwardedFor: r.Header.Get("X-Forwarded-For"),
			XRealIP:       r.Header.Get("X-Real-IP"),
			UserAgent:     r.UserAgent(),
		}

		ctx = context.WithValue(ctx, entity.ClientContextKey{}, clientDevice)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}
