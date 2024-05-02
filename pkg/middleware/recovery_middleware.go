package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"pii-encrypt-example/pkg/exception"
	"pii-encrypt-example/pkg/response"
)

type Recovery struct {
	logger *logrus.Logger
	debug  bool
}

func NewRecovery(logger *logrus.Logger, debug bool) *Recovery {
	return &Recovery{
		logger: logger,
		debug:  debug,
	}
}

func (rm *Recovery) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {

				rm.logger.WithContext(r.Context()).WithFields(
					logrus.Fields{
						"panicking": "true",
					},
				).Error(recovered)

				if rm.debug {
					stack := debug.Stack()
					rm.logger.Error(string(stack))
					w.Header().Set("X-Panic-Response", "true")
				}

				resp := response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "an error occured while processing the request")

				response.JSON(w, resp)
			}

		}()

		handler.ServeHTTP(w, r)
	})
}
