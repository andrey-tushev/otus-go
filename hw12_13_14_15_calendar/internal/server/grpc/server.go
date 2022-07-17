package internalgrpc

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"

	pb "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/protobuf"
)

type Server struct {
	logger     Logger
	app        Application
	grpcServer *grpc.Server
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

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		return errors.New("failed to start grpc listener: " + err.Error())
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterEventsServer(s.grpcServer, pb.UnimplementedEventsServer{})
	err = s.grpcServer.Serve(lsn)
	if err != nil {

		return errors.New("failed to serve rpc: " + err.Error())
	}

	//s.logger.Info("web-server start")
	//
	//s.httpServer = &http.Server{
	//	Addr: net.JoinHostPort(host, port),
	//
	//	ReadTimeout:  1 * time.Second,
	//	WriteTimeout: 1 * time.Second,
	//}
	//s.httpServer.ListenAndServe()
	//
	//<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	s.logger.Info("grpc-server stop")

	return nil
}
