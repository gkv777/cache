package core

import (
	"container/list"

	"github.com/gkv777/cache"
)

// По условиям задания (интерфейс LRUCache) и ключ и значение - строки
type item struct {
	key   string
	value string
}

type LRUBasic struct {
	cap   int
	items map[string]*list.Element // item содержится соотв. в list.Element
	queue *list.List
}

func NewLRUBasicCache(n int) (*LRUBasic, error) {
	if n <= 0 {
		return nil, cache.ErrCapSize
	}
	return &LRUBasic{
		cap:   n,
		items: make(map[string]*list.Element),
		queue: list.New(),
	}, nil
}

func (l *LRUBasic) Add(key, value string) bool {
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

func (l *LRUBasic) Get(key string) (value string, ok bool) {

	e, ok := l.items[key]
	if !ok {
		return "", false
	}

	// Перемещаем связанный с item элемент списка в начало
	l.queue.MoveToFront(e)
	return e.Value.(*item).value, true
}

func (l *LRUBasic) Remove(key string) (ok bool) {

	e, ok := l.items[key]
	if !ok {
		return false
	}

	l.queue.Remove(e)
	delete(l.items, key)
	return true
}

// Вытеснение (удаление) последнего элемента очереди (протухшие данные)
func (l *LRUBasic) removeLast() {
	//l.Lock()
	//defer l.Unlock()

	// получаем последний элемент и удаляем
	if e := l.queue.Back(); e != nil {
		l.queue.Remove(e)
		delete(l.items, e.Value.(*item).key)
	}
}

func (l *LRUBasic) getFirst() *item {
	return l.queue.Front().Value.(*item)
}

func (l *LRUBasic) getLast() *item {
	return l.queue.Back().Value.(*item)
}

func (l *LRUBasic) Cap() int {
	return l.cap
}

func (l *LRUBasic) Len() int {
	return l.queue.Len()
}
