package app

import (
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func eventAppToStorage(event Event) storage.Event {
	return storage.Event{
		ID:       event.ID,
		Title:    event.Title,
		DateTime: event.DateTime,
		Duration: event.Duration,
		Text:     event.Text,
		UserID:   event.UserID,
		Remind:   event.Remind,
	}
}

func eventStorageToApp(event storage.Event) Event {
	return Event{
		ID:       event.ID,
		Title:    event.Title,
		DateTime: event.DateTime,
		Duration: event.Duration,
		Text:     event.Text,
		UserID:   event.UserID,
		Remind:   event.Remind,
	}
}
