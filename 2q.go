package cache

import (
	"sync"
)

const (
	LRU_SIZE = 0.2
	OUT_SIZE = 0.6
)

type TwoQCache struct {
	sync.RWMutex
	fifo *fifo
	lru  *lru
}

func getCacheSizes(n int) (int, int, int) {
	sizeA1out := int(float64(n) * OUT_SIZE)
	sizeAm := int(float64(n) * LRU_SIZE)
	sizeA1in := n - sizeA1out - sizeAm
	return sizeA1in, sizeA1out, sizeAm
}

func NewTwoQCache(n int) (LRUCache, error) {
	if n <= 0 {
		return nil, ErrCapSize
	}

	sizeA1in, sizeA1out, sizeAm := getCacheSizes(n)

	lru, err := newLru(sizeAm)
	if err != nil {
		return nil, err
	}

	fifo, err := newFifo(sizeA1in, sizeA1out)
	if err != nil {
		return nil, err
	}

	return &TwoQCache{		
		fifo: fifo,
		lru: lru,
	}, nil
}

func (c *TwoQCache) Add(key, val string) bool {
	c.Lock()
	defer c.Unlock()
	if c.lru.exists(key) {
		return false
	}
	return c.fifo.Add(key, val)
}

func (c *TwoQCache) Get(key string) (string, bool) {
	c.Lock()
	defer c.Unlock()
	if res, ok := c.lru.Get(key); ok {
		return res, true
	}
	res, ok, out:= c.fifo.Get(key)
	if !ok {
		return "", false
	}
	if out {
		c.lru.Add(key, res)		
	}
	return res, true
}

func (c *TwoQCache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()
	if c.lru.exists(key) {
		return c.lru.Remove(key)
	}
	if c.fifo.exists(key) {
		return c.fifo.Remove(key)
	}
	return false
}

