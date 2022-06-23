package app

import (
	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

func eventAppToStorage(event Event) storage.Event {
	return storage.Event{
		ID:       event.ID,
		Title:    event.Title,
		DateTime: event.DateTime,
		Duration: event.Duration,
		Text:     event.Text,
		UserId:   event.UserId,
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
		UserId:   event.UserId,
		Remind:   event.Remind,
	}
}
