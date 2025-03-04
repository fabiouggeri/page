package util

type void struct{}

var member void

type Set[T comparable] struct {
	items map[T]void
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{items: make(map[T]void, 0)}
	s.AddAll(items...)
	return s
}

func (s *Set[T]) Add(item T) bool {
	_, found := s.items[item]
	if !found {
		s.items[item] = member
		return true
	}
	return false
}

func (s *Set[T]) Contains(items ...T) bool {
	for _, item := range items {
		_, found := s.items[item]
		if !found {
			return false
		}

	}
	return true
}

func (s *Set[T]) AddAll(items ...T) int {
	count := 0
	for _, i := range items {
		if s.Add(i) {
			count++
		}
	}
	return count
}

func (s *Set[T]) Del(item T) bool {
	_, found := s.items[item]
	if found {
		delete(s.items, item)
	}
	return found
}

func (s *Set[T]) Clear() {
	s.items = make(map[T]void, 0)
}

func (s *Set[T]) Items() []T {
	list := make([]T, 0, len(s.items))
	for key := range s.items {
		list = append(list, key)
	}
	return list
}

func (s *Set[T]) Length() int {
	return len(s.items)
}

func (s *Set[T]) Empty() bool {
	return len(s.items) == 0
}

func (s *Set[T]) Equals(other *Set[T]) bool {
	if len(s.items) != len(other.items) {
		return false
	}
	return s.Contains(other.Items()...)
}
