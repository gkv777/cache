package cache

import (
	"sync"
)

const (
	LRU_SIZE = 0.2
	OUT_SIZE = 0.6
)

type TwoQCache struct {
	sync.Mutex
	fifo *fifo
	lru  *lru
}

func NewTwoQCache(n int) (LRUCache, error) {
	if n <= 0 {
		return nil, ErrCapSize
	}

	sizeA1out := int(float64(n) * OUT_SIZE)
	sizeAm := int(float64(n) * LRU_SIZE)
	sizeA1in := n - sizeA1out - sizeAm

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
	return false
}

func (c *TwoQCache) Get(val string) (string, bool) {
	return "", false
}

func (c *TwoQCache) Remove(key string) bool {
	return false
}

func (c *TwoQCache) Cap() int {
	return 0
}

func (c *TwoQCache) Len() int {
	return 0
}
