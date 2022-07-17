package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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
		resp, err := http.Get("http://127.0.0.1:8081/events")
		require.NoError(t, err)

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
