package go_concurrency

import (
	"slices"
	"sync"
)

type ConcurrentSlice[T comparable] struct {
	sync.RWMutex
	items []T
}

type ConcurrentSliceItem[T comparable] struct {
	Index int
	Item  T
}

func NewConcurrentSlice[T comparable](size int) ConcurrentSlice[T] {
	return ConcurrentSlice[T]{
		items: make([]T, size),
	}
}

func (s *ConcurrentSlice[T]) Add(item T) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)
}

func (s *ConcurrentSlice[T]) Get(index int) T {
	s.RLock()
	defer s.RUnlock()

	return s.items[index]
}

func (s *ConcurrentSlice[T]) Has(item T) bool {
	s.Lock()
	defer s.Unlock()

	return slices.Index(s.items, item) >= 0
}

func (s *ConcurrentSlice[T]) RemoveItem(item T) {
	s.Lock()
	defer s.Unlock()
	s.RemoveIndex(slices.Index(s.items, item))
}

func (s *ConcurrentSlice[T]) RemoveIndex(index int) {
	s.Lock()
	defer s.Unlock()
	s.items = append(s.items[:index], s.items[index+1:]...)
}

func (s *ConcurrentSlice[T]) Length() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.items)
}

func (s *ConcurrentSlice[T]) Clear() {
	s.Lock()
	defer s.Unlock()
	s.items = make([]T, 0)
}

func (s *ConcurrentSlice[T]) Iter() <-chan ConcurrentSliceItem[T] {
	iChan := make(chan ConcurrentSliceItem[T], len(s.items))

	f := func() {
		s.RLock()
		defer s.RUnlock()

		for index, item := range s.items {
			iChan <- ConcurrentSliceItem[T]{index, item}
		}
		close(iChan)
	}

	go f()

	return iChan
}

func (s *ConcurrentSlice[T]) All() []T {
	s.RLock()
	defer s.RUnlock()

	return s.items
}
