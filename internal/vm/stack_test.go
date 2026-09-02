package vm

import "testing"

func TestStackLIFO(t *testing.T) {
	stack := newStack[int]()
	stack.push(1)
	stack.push(2)
	if got, err := stack.peek(); err != nil || got != 2 {
		t.Fatalf("peek() = %d, %v", got, err)
	}
	if got, err := stack.pop(); err != nil || got != 2 {
		t.Fatalf("pop() = %d, %v", got, err)
	}
	if got, err := stack.pop(); err != nil || got != 1 {
		t.Fatalf("pop() = %d, %v", got, err)
	}
	if !stack.isEmpty() {
		t.Fatal("stack should be empty")
	}
}

func TestStackEmptyPop(t *testing.T) {
	if _, err := newStack[string]().pop(); err == nil {
		t.Fatal("expected error when popping empty stack")
	}
}
