package hw04_lru_cache //nolint:golint,stylecheck

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
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

	t.Run("purge logic", func(t *testing.T) {
		capacity := 4
		c := NewCache(capacity)

		for i := capacity; i >=0; i-- {
			c.Set(Key(strconv.Itoa(i)), i) //4 [3, 2, 1, 0]
		}
		for i := 0; i < capacity; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}
		val, ok := c.Get(Key(strconv.Itoa(capacity)))
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("most unused", func(t *testing.T) {
		capacity := 4
		c := NewCache(capacity)

		for i := capacity - 1; i >= 0; i-- {
			c.Set(Key(strconv.Itoa(i)), i) //[3, 2, 1, 0]
		}

		for i := 0; i < capacity; i++ {
			c.Get(Key(strconv.Itoa(i))) //[0, 1, 2, 3]
		}

		c.Set(Key(strconv.Itoa(capacity)), capacity)

		val, ok := c.Get(Key(strconv.Itoa(0)))
		require.False(t, ok)
		require.Nil(t, val)

		for i := 1; i <= capacity; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}
	})

	t.Run("clean", func(t *testing.T) {
		capacity := 4
		c := NewCache(capacity)

		keys := make([]Key, 0)
		for i := 0; i < capacity; i++ {
			key := Key(strconv.Itoa(i))
			keys = append(keys, key)
			c.Set(key, i)
		}

		c.Clear()
		for i := 0; i < capacity; i++ {
			val, ok := c.Get(keys[i])
			require.False(t, ok)
			require.Nil(t, val)
		}

		c.Clear()
		for i := 0; i < capacity; i++ {
			val, ok := c.Get(keys[i])
			require.False(t, ok)
			require.Nil(t, val)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
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
