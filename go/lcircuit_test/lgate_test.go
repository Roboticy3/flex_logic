package lcircuit_test

import (
	"container/heap"
	c "flex-logic/collections"
	"flex-logic/lcircuit"
	"testing"
)

func MakeTestPins(states []int) []lcircuit.LPin[int, int] {
	result := []lcircuit.LPin[int, int]{}
	for i := range len(states) {
		result = append(result, lcircuit.LPin[int, int]{
			Nets:  []c.Label{},
			Valid: true,
			State: states[i],
		})
	}
	return result
}

func TestANDGate(t *testing.T) {
	states := MakeTestPins([]int{1, 1, 0}) // Inputs: A=1, B=1, Output=0
	events := &c.LEvents[int, int]{}
	heap.Init(events)

	testGates[0].Solver(states, 0, events) // AND gate

	if states[2].State != 1 {
		t.Errorf("Expected states[2] to be 1, got %d", states[2].State)
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(c.LEvent[int, int])
	if event.Signal != 1 || event.Time != 1 || event.Label != 0 {
		t.Errorf("Unexpected event: %+v", event)
	}
}

func TestNOTGate(t *testing.T) {
	states := MakeTestPins([]int{1, 1, 0}) // Inputs: A=1, Output=0
	events := &c.LEvents[int, int]{}
	heap.Init(events)

	testGates[1].Solver(states, 0, events) // NOT gate

	if states[1].State != ^1 {
		t.Errorf("Expected states[1] to be %d, got %d", ^1, states[1].State)
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(c.LEvent[int, int])
	if event.Signal != ^1 || event.Time != 1 || event.Label != 1 {
		t.Errorf("Unexpected event: %+v", event)
	}
}

func TestLATCHGate(t *testing.T) {
	states := MakeTestPins([]int{1, 0, 0, 0}) // Inputs: SET=1, RESET=0, INNER=0, OUT=0
	events := &c.LEvents[int, int]{}
	heap.Init(events)

	testGates[2].Solver(states, 0, events) // LATCH gate

	if states[2].State != 1 {
		t.Errorf("Expected states[2] to be 1, got %d", states[2].State)
	}
	if states[3].State != 1 {
		t.Errorf("Expected states[3] to be 1, got %d", states[3].State)
	}

	if events.Len() != 1 {
		t.Fatalf("Expected 1 event in the queue, got %d", events.Len())
	}

	event := heap.Pop(events).(c.LEvent[int, int])
	if event.Signal != 1 || event.Time != 2 || event.Label != 2 {
		t.Errorf("Unexpected event: %+v", event)
	}
}
