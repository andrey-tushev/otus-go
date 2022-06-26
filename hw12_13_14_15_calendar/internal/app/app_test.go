package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

func TestCheckAccessibility(t *testing.T) {
	ctx := context.Background()

	logger := &mocks.Logger{}
	logger.On("Info", mock.Anything)

	stor := &mocks.Storage{}
	stor.On("ListEvents", mock.Anything).Return(
		[]storage.Event{
			{
				Title:    "event 1",
				DateTime: time.Date(2022, 12, 25, 0, 0, 0, 0, time.Local),
				Duration: 60 * 60,
			},
			{
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
