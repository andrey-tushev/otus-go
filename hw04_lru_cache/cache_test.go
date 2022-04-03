package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("least recently used", func(t *testing.T) {
		c := NewCache(3)

		require.False(t, c.Set("A", 100))
		require.False(t, c.Set("B", 200))
		require.False(t, c.Set("C", 300))
		// C, B, A

		require.False(t, c.Set("D", 400))
		require.False(t, c.Set("E", 500))
		// E, D, C

		require.True(t, c.Set("C", 301))
		require.True(t, c.Set("D", 401))
		require.True(t, c.Set("E", 501))
		// E, D, C

		require.False(t, c.Set("A", 101))
		require.False(t, c.Set("B", 201))
		// B, A, E

		v, e := c.Get("A")
		require.Equal(t, 101, v)
		require.True(t, e)

		v, e = c.Get("B")
		require.Equal(t, 201, v)
		require.True(t, e)

		v, e = c.Get("C")
		require.Nil(t, v)
		require.False(t, e)

		v, e = c.Get("D")
		require.Nil(t, v)
		require.False(t, e)

		v, e = c.Get("E")
		require.Equal(t, 501, v)
		require.True(t, e)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(10)

		c.Set("A", 100)
		c.Set("B", 200)
		c.Set("C", 300)
		c.Clear()

		v, e := c.Get("B")
		require.Nil(t, v)
		require.False(t, e)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
