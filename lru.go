package cache

import "sync"

type LRU struct {
	sync.RWMutex
	lru *lru
}

func NewLRUCache(n int) (LRUCache, error) {
	basic, err := newLru(n)
	if err != nil {
		return nil, err
	}
	return &LRU{
		lru:       basic,
	}, nil
}

func (l *LRU) Add(key, value string) bool {
	l.Lock()
	defer l.Unlock()

	return l.lru.Add(key, value)
}

func (l *LRU) Get(key string) (value string, ok bool) {
	l.Lock()
	defer l.Unlock()
	return l.lru.Get(key)
}

func (l *LRU) Remove(key string) (ok bool) {
	l.Lock()
	defer l.Unlock()
	return l.lru.Remove(key)
}

