package cache

import (
	"sync"
)

type TwoQCache struct {
	sync.Mutex
	qIn  LRUBasic
	qOut LRUBasic
	qHot LRUBasic
}
