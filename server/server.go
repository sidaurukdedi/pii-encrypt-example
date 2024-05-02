package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	startingMessage string = "HTTP Server starts to listen on %s"
	shutdownMessage string = "HTTP Server is gracefully shutdown."
)

// Server is a concrete struct of http server.
type Server struct {
	logger     *logrus.Logger
	httpServer *http.Server
}

// NewServer is a constructor.
func NewServer(logger *logrus.Logger, handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	return &Server{
		logger:     logger,
		httpServer: httpServer,
	}
}

// Start will start the server.
// Do not call this in goroutine.
func (s *Server) Start() {
	go func() {
		s.logger.Info(fmt.Sprintf(startingMessage, s.httpServer.Addr))
		s.httpServer.ListenAndServe()
	}()
}

// Close will block all the incomming request and subsequently shutdown the server.
func (s *Server) Close() {
	s.httpServer.Shutdown(context.Background())
	s.logger.Info(shutdownMessage)
}
