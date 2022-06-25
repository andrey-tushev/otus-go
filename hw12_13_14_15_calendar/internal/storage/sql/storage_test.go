package sqlstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	stor := New("postgres://calendar:calendar@localhost/calendar")
	err := stor.Connect(ctx)
	require.NoError(t, err)
	err = stor.Exec(ctx, "DELETE FROM events")
	require.NoError(t, err)

	list, err := stor.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, list, 0)

	id1, err := stor.CreateEvent(ctx, storage.Event{Title: "Title 1"})
	require.NoError(t, err)
	require.NotEmpty(t, id1)

	id2, err := stor.CreateEvent(ctx, storage.Event{Title: "Title 2"})
	require.NoError(t, err)
	require.NotEmpty(t, id2)

	id3, err := stor.CreateEvent(ctx, storage.Event{Title: "Title 3"})
	require.NoError(t, err)
	require.NotEmpty(t, id3)

	list, err = stor.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, list, 3)

	err = stor.UpdateEvent(ctx, storage.Event{ID: id2, Title: "Title 2A"})
	require.NoError(t, err)

	err = stor.DeleteEvent(ctx, id1)
	require.NoError(t, err)

	list, err = stor.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)

	for _, item := range list {
		if item.ID == id2 {
			require.Equal(t, "Title 2A", item.Title)
		} else if item.ID == id3 {
			require.Equal(t, "Title 3", item.Title)
		} else {
			t.Fatal("bad id in the list")
		}
	}

	err = stor.Exec(ctx, "DELETE FROM events")
	require.NoError(t, err)
}