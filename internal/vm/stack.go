package vm

import "errors"

type stack[T any] struct {
	items []T
}

func newStack[T any]() *stack[T] {
	return &stack[T]{items: []T{}}
}

func (s *stack[T]) push(value T) {
	s.items = append(s.items, value)
}

func (s *stack[T]) pop() (T, error) {
	if s.isEmpty() {
		var zeroValue T
		return zeroValue, errors.New("stack is empty")
	}
	last := len(s.items) - 1
	value := s.items[last]
	s.items = s.items[:last]
	return value, nil
}

func (s *stack[T]) peek() (T, error) {
	if s.isEmpty() {
		var zeroValue T
		return zeroValue, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *stack[T]) isEmpty() bool {
	return len(s.items) == 0
}
