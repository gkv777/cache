package cache

import (
	"container/list"

)

type item struct {
	key   string
	value string
}

type lru struct {
	cap   int
	items map[string]*list.Element 
	queue *list.List
}

func newLru(n int) (*lru, error) {
	if n <= 0 {
		return nil, ErrCapSize
	}
	return &lru{
		cap:   n,
		items: make(map[string]*list.Element),
		queue: list.New(),
	}, nil
}


func (l *lru) Add(key, value string) bool {
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

	e, ok := l.items[key]
	if !ok {
		return "", false
	}

	// Перемещаем связанный с item элемент списка в начало
	l.queue.MoveToFront(e)
	return e.Value.(*item).value, true
}

func (l *lru) Remove(key string) (ok bool) {

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

func (l *lru) exists(key string) bool {
	_, ok := l.items[key]
	return ok
}

func (l *lru) getFirst() *item {
	return l.queue.Front().Value.(*item)
}

func (l *lru) getLast() *item {
	return l.queue.Back().Value.(*item)
}

func (l *lru) Cap() int {
	return l.cap
}

func (l *lru) Len() int {
	return l.queue.Len()
}
