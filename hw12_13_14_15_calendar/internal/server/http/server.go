package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	logger     Logger
	app        Application
	httpServer *http.Server
}

type Logger interface {
	Info(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) (string, error)
	UpdateEvent(ctx context.Context, event app.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]app.Event, error)
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
	router := chi.NewRouter()

	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello-world"))
	})
	router.Route("/events", func(subRouter chi.Router) {
		subRouter.Get("/", s.ListEvents)
		subRouter.Post("/", s.CreateEvent)
		subRouter.Patch("/", s.UpdateEvent)
	})

	router.ServeHTTP(w, r)
}
