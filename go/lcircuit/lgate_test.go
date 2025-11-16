package lcircuit

import (
	"container/heap"
	"testing"
)

func TestANDGate(t *testing.T) {
	states := []int{1, 1, 0} // Inputs: A=1, B=1, Output=0
	events := &levents[int, int]{}
	heap.Init(events)

	testGates[0].Solver(states, 0, events) // AND gate

	if states[2] != 1 {
		t.Errorf("Expected states[2] to be 1, got %d", states[2])
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(levent[int, int])
	if event.signal != 1 || event.time != 1 || event.label != 0 {
		t.Errorf("Unexpected event: %+v", event)
	}
}

func TestNOTGate(t *testing.T) {
	states := []int{1, 0} // Input: A=1, Output=0
	events := &levents[int, int]{}
	heap.Init(events)

	testGates[1].Solver(states, 0, events) // NOT gate

	if states[1] != ^1 {
		t.Errorf("Expected states[1] to be %d, got %d", ^1, states[1])
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(levent[int, int])
	if event.signal != ^1 || event.time != 1 || event.label != 1 {
		t.Errorf("Unexpected event: %+v", event)
	}
}

func TestLATCHGate(t *testing.T) {
	states := []int{1, 0, 0, 0} // Inputs: SET=1, RESET=0, INNER=0, OUT=0
	events := &levents[int, int]{}
	heap.Init(events)

	testGates[2].Solver(states, 0, events) // LATCH gate

	if states[2] != 1 {
		t.Errorf("Expected states[2] to be 1, got %d", states[2])
	}
	if states[3] != 1 {
		t.Errorf("Expected states[3] to be 1, got %d", states[3])
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(levent[int, int])
	if event.signal != 1 || event.time != 2 || event.label != 2 {
		t.Errorf("Unexpected event: %+v", event)
	}
}
