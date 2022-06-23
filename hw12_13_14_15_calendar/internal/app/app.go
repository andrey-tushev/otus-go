package app

import (
	"context"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface { // TODO
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (string, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) (string, error) {
	return a.storage.CreateEvent(ctx, eventAppToStorage(event))
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	return a.storage.UpdateEvent(ctx, eventAppToStorage(event))
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context) ([]Event, error) {
	list, err := a.storage.ListEvents(ctx)
	if err != nil {
		return []Event{}, nil
	}

	var events = make([]Event, 0, len(list))
	for _, e := range list {
		events = append(events, eventStorageToApp(e))
	}
	return events, nil
}
