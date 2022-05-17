package cache

import "container/list"

type fifo struct {
	inCap  int
	outCap int
	in     map[string]*list.Element
	out    map[string]*list.Element
	qIn    *list.List
	qOut   *list.List
}

func newFifo(in, out int) (*fifo, error) {
	if in <= 0 || out <= 0 {
		return nil, ErrCapSize
	}
	return &fifo{
		inCap:  in,
		outCap: out,
		in:     make(map[string]*list.Element),
		out:    make(map[string]*list.Element),
		qIn:    list.New(),
		qOut:   list.New(),
	}, nil
}

func (f *fifo) exists(key string) bool {
	if _, ok := f.in[key]; ok {
		return true
	}
	_, ok := f.out[key]
	return ok
}

// Добавление к A1in:
//
func (f *fifo) Add(key, value string) bool {
	if f.exists(key) {
		return false
	}
	// Если А1in заполнен - вытесняем в A1out
	if f.inLen() == f.inCap {
		f.moveLastToOut()
	}
	// Добавляем в начало списка новые данне
	i := &item{
		key:   key,
		value: value,
	}
	e := f.qIn.PushFront(i)

	f.in[key] = e
	return true
}

func (f *fifo) Get(key string) (string, bool, bool) {
	if e, ok := f.in[key]; ok {
		return e.Value.(*item).value, true, false
	}
	if e, ok:= f.out[key]; ok {
		// удаляем из очереди (для переноса в Am)
		f.qOut.Remove(e)
		delete(f.out, key)
		return e.Value.(*item).value, true, true
	}
	return "", false, false
}

func (f *fifo)Remove(key string) bool {
	if e, ok := f.in[key]; ok {
		f.qIn.Remove(e)
		delete(f.in, key)
		return true
	}
	if e, ok := f.out[key]; ok {
		f.qOut.Remove(e)
		delete(f.out, key)
		return true
	}
	return false
}

// Вытеснение (удаление) последнего элемента очереди (протухшие данные)
func (f *fifo) moveLastToOut() {
	// если out заполнен - удалям последнюю запись
	if f.outLen() == f.outCap {
		if l := f.qOut.Back(); l != nil {
			f.qOut.Remove(l)
			delete(f.out, l.Value.(*item).key)
		}
	}
	// получаем последний элемент в in, удаляем
	if e := f.qIn.Back(); e != nil {
		f.qIn.Remove(e)
		delete(f.in, e.Value.(*item).key)
		// добавляем в out
		i:= f.qOut.PushFront(e.Value.(*item))
		f.out[e.Value.(*item).key] = i
	}
}

func (f *fifo) inLen() int {
	return f.qIn.Len()
}

func (f *fifo) outLen() int {
	return f.qOut.Len()
}

func (f *fifo) Len() int {
	return f.qIn.Len() + f.qOut.Len()
}

func(f *fifo) Cap() int {
	return f.inCap + f.outCap
}