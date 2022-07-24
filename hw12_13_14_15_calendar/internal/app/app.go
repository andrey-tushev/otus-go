package app

import (
	"context"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (string, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]storage.Event, error)
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) (string, error) {
	if err := a.checkAccessibility(ctx, event); err != nil {
		return "", err
	}
	return a.storage.CreateEvent(ctx, eventAppToStorage(event))
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	if err := a.checkAccessibility(ctx, event); err != nil {
		return err
	}
	return a.storage.UpdateEvent(ctx, eventAppToStorage(event))
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context) ([]Event, error) {
	list, err := a.storage.ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(list))
	for _, e := range list {
		events = append(events, eventStorageToApp(e))
	}
	return events, nil
}

func (a *App) checkAccessibility(ctx context.Context, event Event) error {
	list, err := a.storage.ListEvents(ctx)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.ID == event.ID {
			continue
		}

		if event.DateTime.Unix()+int64(event.Duration) >= item.DateTime.Unix() &&
			event.DateTime.Unix() <= item.DateTime.Unix()+int64(item.Duration) {
			return ErrDateBusy
		}
	}

	return nil
}

func (a *App) Close(ctx context.Context) {
	a.storage.Close(ctx)
}
