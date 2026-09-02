package utils

import (
	"errors"
)

// Stack represents a stack data structure.
type Stack[T any] struct {
	items []T
}

// Constructor
func New[T any]() *Stack[T] {
	return &Stack[T]{items: []T{}}
}

// Push adds an item to the stack.
func (s *Stack[T]) Push(data T) error {
	s.items = append(s.items, data)
	return nil
}

// Pop removes and returns the top item from the stack.
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New("pop failed: stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Peek looks ahead without removing from the stack
func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New("peek failed: stack is empty")
	}
	item := s.items[len(s.items)-1]
	return item, nil
}

// IsEmpty checks if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}
