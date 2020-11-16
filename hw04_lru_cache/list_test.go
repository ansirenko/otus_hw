package hw04_lru_cache //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	seedRandom()

	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			val := i.Value.(int)
			elems = append(elems, val)
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("delete last", func(t *testing.T) {
		l := NewList()
		elem := rand.Int()
		listItem := l.PushFront(elem)
		l.Remove(listItem)

		require.Zero(t, l.Len())
		require.Nil(t, l.Back())
		require.Nil(t, l.Front())
	})

	t.Run("move single element", func(t *testing.T) {
		l := NewList()
		elem := rand.Float64()
		listItem := l.PushFront(elem)
		l.MoveToFront(listItem)

		require.Equal(t, listItem, l.Front())
		require.Equal(t, listItem, l.Back())
	})

	t.Run("len", func(t *testing.T) {
		l := NewList()

		require.Zero(t, l.Len())

		count := rand.Int31n(200)
		for i := 0; i < int(count); i++ {
			l.PushFront(i)
		}
		require.EqualValues(t, count, l.Len())

		for i := 0; i < int(count); i++ {
			l.PushBack(i)
		}
		require.EqualValues(t, count * 2, l.Len())

		for i := 0; i < int(count * 2); i++ {
			l.MoveToFront(l.Back())
		}
		require.EqualValues(t, count * 2, l.Len())

		for i := 0; i < int(count * 2); i++ {
			l.Remove(l.Front())
		}
		require.Zero(t, l.Len())
	})

	t.Run("run", func(t *testing.T) {
		l := NewList()

		val := l.Front()
		require.Nil(t, val)

		item := l.PushBack(0)
		val = l.Front()
		require.Equal(t, item, val)

		item2 := l.PushFront(1)
		val = l.Front()
		require.Equal(t, item2, val)

		l.MoveToFront(item)
		val = l.Front()
		require.Equal(t, item, val)

		l.Remove(item)
		val = l.Front()
		require.Equal(t, item2, val)
	})

	t.Run("back", func(t *testing.T) {
		l := NewList()

		val := l.Back()
		require.Nil(t, val)

		item := l.PushBack(0)
		val = l.Back()
		require.Equal(t, item, val)

		l.Remove(item)
		val = l.Back()
		require.Nil(t, val)

		item = l.PushFront(0)
		val = l.Back()
		require.Equal(t, item, val)

		item2 := l.PushBack(1)
		val = l.Back()
		require.Equal(t, item2, val)

		l.MoveToFront(item2)
		val = l.Back()
		require.Equal(t, item, val)
	})

	t.Run("push front", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(0)
		val := l.Front()
		require.Equal(t, item, val)
		val = l.Back()
		require.Equal(t, item, val)
		length := l.Len()
		require.Equal(t, 1, length)

		item2 := l.PushFront(1)
		val = l.Front()
		require.Equal(t, item2, val)
		next := item2.Next
		require.Equal(t, item, next)
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()

		item := l.PushBack(0)
		val := l.Back()
		require.Equal(t, item, val)
		val = l.Front()
		require.Equal(t, item, val)
		length := l.Len()
		require.Equal(t, 1, length)

		item2 := l.PushBack(1)
		val = l.Back()
		require.Equal(t, item2, val)
		prev := item2.Prev
		require.Equal(t, item, prev)
	})

	t.Run("remove", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(0)
		l.Remove(item)
		length := l.Len()
		require.Zero(t, length)

		item = l.PushBack(1)
		l.Remove(item)
		length = l.Len()
		require.Zero(t, length)

		items := make(map[int]*ListItem, 0)
		for i := 0; i < 3; i++ {
			items[i] = l.PushFront(i)
		}
		l.Remove(items[1])
		length = l.Len()
		require.Equal(t, len(items) - 1, length)
		prev := items[0].Prev
		require.Equal(t, items[2], prev)
		next := items[2].Next
		require.Equal(t, items[0], next)
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(0)
		l.MoveToFront(item)
		val := l.Front()
		require.Equal(t, item, val)
		val = l.Back()
		require.Equal(t, item, val)

		item2 := l.PushBack(1)
		l.MoveToFront(item2)
		val = l.Front()
		require.Equal(t, item2, val)

		l.PushBack(2)
		l.MoveToFront(item)
		val = l.Front()
		require.Equal(t, item, val)
	})
}

func seedRandom() {
	seed := time.Now().Unix()
	fmt.Printf("seed: %d\n", seed)
	rand.Seed(seed)
}
