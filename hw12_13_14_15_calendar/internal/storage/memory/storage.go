package memorystorage

import (
	"context"
	"sync"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/app"
	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"

	"github.com/google/uuid"
)

type Storage struct {
	mu     *sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     &sync.RWMutex{},
		events: make(map[string]storage.Event),
	}
}

func (s Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	event.ID = id
	s.events[id] = event
	return id, nil
}

func (s Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, has := s.events[event.ID]
	if !has {
		return app.ErrNotFound
	}
	s.events[event.ID] = event

	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, has := s.events[id]
	if !has {
		return app.ErrNotFound
	}
	delete(s.events, id)

	return nil
}

func (s Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		list = append(list, event)
	}

	return list, nil
}
