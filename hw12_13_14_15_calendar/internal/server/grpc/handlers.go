package internalgrpc

import (
	"context"
	"time"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	pb "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/protobuf"
)

func (s *Server) CreateEvent(ctx context.Context, request *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	id, err := s.app.CreateEvent(ctx, app.Event{
		ID:       request.Event.ID,
		Title:    request.Event.Title,
		DateTime: time.Unix(request.Event.DateTime, 0),
		Duration: int(request.Event.Duration),
		Text:     request.Event.Text,
		UserID:   int(request.Event.UserID),
		Remind:   int(request.Event.Remind),
	})

	s.logger.Info("grpc CreateEvent")

	return &pb.CreateEventResponse{ID: id}, err
}

func (s *Server) UpdateEvent(ctx context.Context, request *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	err := s.app.UpdateEvent(ctx, app.Event{
		ID:       request.Event.ID,
		Title:    request.Event.Title,
		DateTime: time.Unix(request.Event.DateTime, 0),
		Duration: int(request.Event.Duration),
		Text:     request.Event.Text,
		UserID:   int(request.Event.UserID),
		Remind:   int(request.Event.Remind),
	})

	s.logger.Info("grpc UpdateEvent")

	return &pb.UpdateEventResponse{}, err
}

func (s *Server) ListEvents(ctx context.Context, request *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	list, err := s.app.ListEvents(ctx)

	events := make([]*pb.Event, 0, len(list))
	for _, item := range list {
		events = append(events, &pb.Event{
			ID:       item.ID,
			Title:    item.Title,
			DateTime: item.DateTime.Unix(),
			Duration: int64(item.Duration),
			Text:     item.Text,
			UserID:   int64(item.UserID),
			Remind:   int64(item.Remind),
		})
	}

	s.logger.Info("grpc ListEvents")

	return &pb.ListEventsResponse{Events: events}, err
}
