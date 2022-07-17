package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	logger     Logger
	app        Application
	httpServer *http.Server
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

func (s *Server) Start(ctx context.Context, host string, port string) error {
	s.logger.Info("web-server start")

	s.httpServer = &http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      s.loggingMiddleware(s),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.httpServer.ListenAndServe()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.Shutdown(ctx)
	s.logger.Info("web-server stopped")

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
