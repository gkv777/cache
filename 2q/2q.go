package two_q

import (
	"sync"

	"github.com/gkv777/cache/core"
)

type TwoQCache struct {
	sync.Mutex
	qIn core.LRUBasic
	qOut core.LRUBasic
	qHot core.LRUBasic
}