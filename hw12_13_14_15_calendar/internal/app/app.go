package app

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/queue/message"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

const cleanAge = 365 * 24 * time.Hour

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
	// ListEvents хорошо бы научить возвращать список за определенный период, но для учебной задачи оставим как есть
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

func (a *App) Clean(ctx context.Context) {
	a.logger.Info("Cleanup")

	events, err := a.ListEvents(ctx)
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, event := range events {
		since := time.Since(event.DateTime)
		if since > cleanAge {
			a.logger.Info(fmt.Sprintf("Cleaning: %s, %s", event.DateTime, event.Title))

			err := a.DeleteEvent(ctx, event.ID)
			if err != nil {
				a.logger.Error(err.Error())
				return
			}
		}
	}
}

func (a *App) Remind(ctx context.Context, producer *rabbitmq.Producer, interval time.Duration) {
	a.logger.Info("Remind")

	events, err := a.ListEvents(ctx)
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, event := range events {
		// Лучше было бы помечать в сторадже, что напоминание отправлено
		// Но для учебных целей оставим так
		if needRemind(event.DateTime, interval) {
			a.logger.Info(fmt.Sprintf("Remind: %s, %s", event.DateTime, event.Title))
		}

		message := message.New(event.Title, event.DateTime)
		producer.Publish(message)
	}
}

func needRemind(t time.Time, interval time.Duration) bool {
	now := time.Now()
	if now.Unix() < t.Unix() && t.Unix() < now.Add(interval).Unix() {
		return true
	}
	return false
}

func (a *App) Close(ctx context.Context) {
	a.storage.Close(ctx)
}
