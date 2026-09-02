package utils

import "testing"

func TestStackLIFO(t *testing.T) {
	stack := New[int]()
	if err := stack.Push(1); err != nil {
		t.Fatal(err)
	}
	if err := stack.Push(2); err != nil {
		t.Fatal(err)
	}
	if got, err := stack.Peek(); err != nil || got != 2 {
		t.Fatalf("Peek() = %d, %v", got, err)
	}
	if got, err := stack.Pop(); err != nil || got != 2 {
		t.Fatalf("Pop() = %d, %v", got, err)
	}
	if got, err := stack.Pop(); err != nil || got != 1 {
		t.Fatalf("Pop() = %d, %v", got, err)
	}
	if !stack.IsEmpty() {
		t.Fatal("stack should be empty")
	}
}

func TestStackEmptyPop(t *testing.T) {
	if _, err := New[string]().Pop(); err == nil {
		t.Fatal("expected error when popping empty stack")
	}
}
