package internalhttp

// nolint:gci
import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/server/http/mocks"
)

func TestList(t *testing.T) {
	// Запустим сервер с замоканным Application
	{
		application := &mocks.Application{}
		application.On("ListEvents", mock.Anything).Return([]app.Event{{
			ID:       "123",
			Title:    "The title",
			DateTime: time.Time{},
			Duration: 600,
			Text:     "The text",
			UserID:   17,
			Remind:   300,
		}}, nil)

		server := NewServer(&logger.Logger{}, application)
		ctx, cancel := context.WithCancel(context.Background())
		go server.Start(ctx, "127.0.0.1", "8081")

		defer cancel()
	}

	// Обратимся к серверу, прочитаем список
	{
		// nolint:noctx
		resp, err := http.Get("http://127.0.0.1:8081/events")
		require.NoError(t, err)
		defer resp.Body.Close()

		responseJSON, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var events []app.Event
		err = json.Unmarshal(responseJSON, &events)
		require.NoError(t, err)

		require.Len(t, events, 1)
		require.Equal(t, "123", events[0].ID)
		require.Equal(t, "The title", events[0].Title)
		require.Equal(t, "The text", events[0].Text)
	}
}

func TestCreate(t *testing.T) {
	application := &mocks.Application{}

	testTime, err := time.Parse(time.RFC3339, "2018-03-02T15:04:05Z")
	require.NoError(t, err)

	application.On(
		"CreateEvent",
		mock.Anything,
		app.Event{
			Title:    "The title",
			DateTime: testTime,
			Duration: 600,
			Text:     "A text description",
			UserID:   100,
			Remind:   300,
		},
	).Return("new-id", nil)

	// Запустим сервер с замоканным Application
	{
		server := NewServer(&logger.Logger{}, application)
		ctx, cancel := context.WithCancel(context.Background())
		go server.Start(ctx, "127.0.0.1", "8081")

		defer cancel()
	}

	// Отправим в API команду на создание записи
	{
		// nolint:noctx
		resp, err := http.Post(
			"http://127.0.0.1:8081/events",
			"application/json",
			strings.NewReader(`{
				"Title":    "The title",
				"DateTime": "2018-03-02T15:04:05Z",
				"Duration": 600,
				"Text":     "A text description",
				"UserID":   100, 
				"Remind":   300
			}`))
		require.NoError(t, err)
		defer resp.Body.Close()

		responseJSON, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var event app.Event
		err = json.Unmarshal(responseJSON, &event)
		require.NoError(t, err)

		require.Equal(t, "new-id", event.ID)
		require.Equal(t, "The title", event.Title)
		require.Equal(t, "A text description", event.Text)
	}
}
