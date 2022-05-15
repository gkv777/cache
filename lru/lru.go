package lru

import (
	"sync"

	"github.com/gkv777/cache"
	"github.com/gkv777/cache/core"
)

// По условиям задания (интерфейс LRUCache) и ключ и значение - строки
type item struct {
	key   string
	value string
}

type LRU struct {
	sync.RWMutex
	q *core.LRUBasic
}

func NewLRUCache(n int) (cache.LRUCache, error) {
	basic, err := core.NewLRUBasicCache(n)
	if err != nil {
		return nil, err
	}
	return &LRU{
		RWMutex: sync.RWMutex{},
		q:       basic,
	}, nil
}

func (l *LRU) Add(key, value string) bool {
	l.Lock()
	defer l.Unlock()

	return l.q.Add(key, value)
}

func (l *LRU) Get(key string) (value string, ok bool) {
	l.Lock()
	defer l.Unlock()
	return l.q.Get(key)
}

func (l *LRU) Remove(key string) (ok bool) {
	l.Lock()
	defer l.Unlock()
	return l.q.Remove(key)
}

func (l *LRU) Cap() int {
	l.RLock()
	defer l.RUnlock()
	return l.q.Cap()
}

func (l *LRU) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.q.Len()
}
