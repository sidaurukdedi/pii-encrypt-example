package entity

type ClientContextKey struct{}

type ClientDevice struct {
	RemoteAddress string
	XForwardedFor string
	XRealIP       string
	UserAgent     string
}
