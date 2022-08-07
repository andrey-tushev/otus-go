package internalgrpc

// nolint
import (
	"context"
	"errors"
	"net"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	pb "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/protobuf"

	"google.golang.org/grpc"
)

type Server struct {
	logger     Logger
	app        Application
	grpcServer *grpc.Server

	pb.UnimplementedEventsServer
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

func (s *Server) Start(ctx context.Context, port string) error {
	lsn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errors.New("failed to start grpc listener: " + err.Error())
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterEventsServer(s.grpcServer, s)
	err = s.grpcServer.Serve(lsn)
	if err != nil {
		return errors.New("failed to serve rpc: " + err.Error())
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	s.logger.Info("grpc-server stop")

	return nil
}
