package memorystorage

//nolint:gci
import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	stor := New()

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
		switch item.ID {
		case id2:
			require.Equal(t, "Title 2A", item.Title)
		case id3:
			require.Equal(t, "Title 3", item.Title)
		default:
			t.Fatal("bad id in the list")
		}
	}
}

func TestMultithreading(t *testing.T) {
	ctx := context.Background()
	stor := New()

	const threads = 10
	const records = 1000

	wg := &sync.WaitGroup{}
	wg.Add(threads)

	for t := 0; t < threads; t++ {
		go func() {
			defer wg.Done()
			for i := 0; i < records; i++ {
				stor.CreateEvent(ctx, storage.Event{Title: strconv.Itoa(i)})
			}
		}()
	}

	wg.Wait()

	list, err := stor.ListEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, threads*records, len(list))
}
