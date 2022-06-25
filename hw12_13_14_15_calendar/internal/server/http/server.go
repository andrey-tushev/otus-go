package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	logger Logger
	app    Application
}

type Logger interface {
	Info(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("web-server start")

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      s.loggingMiddleware(s),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	server.ListenAndServe()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("web-server stop")

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/hello":
		w.Write([]byte("hello-world"))

	default:
		http.NotFound(w, r)

	}
}
