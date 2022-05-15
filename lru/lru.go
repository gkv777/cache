package lru

import (
	"container/list"
	"sync"

	"github.com/gkv777/cache"
)

// По условиям задания (интерфейс LRUCache) и ключ и значение - строки
type item struct {
	key   string
	value string
}

type lru struct {
	sync.RWMutex
	cap   int
	items map[string]*list.Element // item содержится соотв. в list.Element
	queue *list.List
}

func NewLRUCache(n int) (cache.LRUCache, error) {
	if n <= 0 {
		return nil, cache.ErrCapSize
	}
	return &lru{
		cap:   n,
		items: make(map[string]*list.Element),
		queue: list.New(),
	}, nil
}

func (l *lru) Add(key, value string) bool {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.items[key]; ok {
		return false
	}
	// Если
	if l.queue.Len() == l.cap {
		l.removeLast()
	}
	// Добавляем в начало списка новые данне
	i := &item{
		key:   key,
		value: value,
	}
	e := l.queue.PushFront(i)

	l.items[key] = e
	return true
}

func (l *lru) Get(key string) (value string, ok bool) {
	l.Lock()
	defer l.Unlock()

	e, ok := l.items[key]
	if !ok {
		return "", false
	}

	// Перемещаем связанный с item элемент списка в начало
	l.queue.MoveToFront(e)
	return e.Value.(*item).value, true
}

func (l *lru) Remove(key string) (ok bool) {
	l.Lock()
	defer l.Unlock()

	e, ok := l.items[key]
	if !ok {
		return false
	}

	l.queue.Remove(e)
	delete(l.items, key)
	return true
}

// Вытеснение (удаление) последнего элемента очереди (протухшие данные)
func (l *lru) removeLast() {
	//l.Lock()
	//defer l.Unlock()

	// получаем последний элемент и удаляем
	if e := l.queue.Back(); e != nil {
		l.queue.Remove(e)
		delete(l.items, e.Value.(*item).key)
	}
}

func (l *lru) getFirst() *item {
	l.RLock()
	defer l.RUnlock()
	return l.queue.Front().Value.(*item)
}

func (l *lru) getLast() *item {
	l.RLock()
	defer l.RUnlock()
	return l.queue.Back().Value.(*item)
}

func (l *lru) Cap() int {
	l.RLock()
	defer l.RUnlock()
	return l.cap
}

func (l *lru) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.queue.Len()
}
