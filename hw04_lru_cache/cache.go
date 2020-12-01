package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	item, ok := lru.items[key]
	if ok {
		item.Value = cacheItem{
			key:   key,
			value: value,
		}
		lru.queue.PushFront(item)
		lru.items[key] = item
		return true
	}
	if lru.capacity == len(lru.items) {
		lastItem := lru.queue.Back()
		lru.queue.Remove(lastItem)
		keyToDelete := lastItem.Value.(cacheItem).key
		delete(lru.items, keyToDelete)
	}
	item = lru.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	lru.items[key] = item
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	item, ok := lru.items[key]
	if !ok {
		return nil, false
	}
	lru.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, true
}

func (lru *lruCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.items = make(map[Key]*ListItem)
	lru.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}
