package app

//nolint:gci
import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func TestCheckAccessibility(t *testing.T) {
	ctx := context.Background()

	logger := &mocks.Logger{}
	logger.On("Info", mock.Anything)

	stor := &mocks.Storage{}
	stor.On("ListEvents", mock.Anything).Return(
		[]storage.Event{
			{
				ID:       "1001",
				Title:    "event 1",
				DateTime: time.Date(2022, 12, 25, 0, 0, 0, 0, time.Local),
				Duration: 60 * 60,
			},
			{
				ID:       "1002",
				Title:    "event 2",
				DateTime: time.Date(2022, 12, 26, 0, 0, 0, 0, time.Local),
				Duration: 60 * 60,
			},
		}, nil)

	app := New(logger, stor)

	err := app.checkAccessibility(ctx, Event{
		Title:    "new event",
		DateTime: time.Date(2022, 12, 25, 2, 0, 0, 0, time.Local),
		Duration: 60 * 60,
	})
	require.NoError(t, err)

	err = app.checkAccessibility(ctx, Event{
		Title:    "new event",
		DateTime: time.Date(2022, 12, 25, 0, 30, 0, 0, time.Local),
		Duration: 60 * 60,
	})
	require.Error(t, err)
}

func TestNeedRemind(t *testing.T) {
	eventTime := time.Now().Add(30 * time.Minute)
	require.True(t, needRemind(eventTime, 1*time.Hour))

	eventTime = time.Now().Add(65 * time.Minute)
	require.False(t, needRemind(eventTime, 1*time.Hour))

	eventTime = time.Now().Add(-5 * time.Minute)
	require.False(t, needRemind(eventTime, 1*time.Hour))
}
