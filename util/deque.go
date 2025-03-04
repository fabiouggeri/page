package util

import "errors"

type Deque[T any] struct {
	items []T
	empty T
}

func NewDeque[T any](items ...T) *Deque[T] {
	var capacity int
	if len(items) > 0 {
		capacity = int(float64(len(items)) * float64(1.3))
	} else {
		capacity = 8
	}
	d := &Deque[T]{items: make([]T, 0, capacity)}
	d.items = append(d.items, items...)
	return d
}

func (q *Deque[T]) Add(items ...T) {
	q.AddLast(items...)
}

func (q *Deque[T]) AddFirst(items ...T) {
	q.items = append(items, q.items...)
}

func (q *Deque[T]) AddLast(items ...T) {
	q.items = append(q.items, items...)
}

func (q *Deque[T]) Push(items ...T) {
	q.items = append(q.items, items...)
}

func (q *Deque[T]) Pop() (T, error) {
	return q.Last()
}

func (q *Deque[T]) First() (T, error) {
	if len(q.items) == 0 {
		return q.empty, errors.New("empty stack")
	}
	item := q.items[0]
	q.items[0] = q.empty
	q.items = q.items[1:]
	return item, nil
}

func (q *Deque[T]) PeekFirst() (T, error) {
	if len(q.items) == 0 {
		return q.empty, errors.New("empty stack")
	}
	return q.items[0], nil
}

func (q *Deque[T]) Last() (T, error) {
	if len(q.items) == 0 {
		return q.empty, errors.New("empty stack")
	}
	item := q.items[len(q.items)-1]
	q.items[len(q.items)-1] = q.empty
	q.items = q.items[:len(q.items)-1]
	return item, nil
}

func (q *Deque[T]) PeekLast() (T, error) {
	if len(q.items) == 0 {
		return q.empty, errors.New("empty stack")
	}
	return q.items[len(q.items)-1], nil
}

func (q *Deque[T]) Clear() {
	q.items = make([]T, 0)
}

func (q *Deque[T]) Empty() bool {
	return len(q.items) == 0
}

func (q *Deque[T]) Len() int {
	return len(q.items)
}
