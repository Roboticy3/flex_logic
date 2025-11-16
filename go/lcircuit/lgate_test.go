package Lcircuit

import (
	"container/heap"
	"testing"
)

var names Llabeling[string_component] = Llabeling[string_component]{
	"AND",
	"NOT",
	"LATCH",
}

var testGates []Lgate[int, int] = []Lgate[int, int]{
	{
		name: 0,
		Solver: func(states []int, time int, events *levents[int, int]) {
			states[2] = states[0] & states[1]
			events.Push(levent[int, int]{
				signal: states[2],
				time:   time + 1,
				label:  0,
			})
		},
		pinout: []string{"A", "B", "OUT"},
	},
	{
		name: 1,
		Solver: func(states []int, time int, events *levents[int, int]) {
			states[1] = ^states[0]
			events.Push(levent[int, int]{
				signal: states[1],
				time:   time + 1,
				label:  1,
			})
		},
		pinout: []string{"A", "B", "OUT"},
	},
	{
		name: 2,
		Solver: func(states []int, time int, events *levents[int, int]) {
			states[2] |= states[0]
			states[2] &^= states[1]
			states[3] = states[2]
			events.Push(levent[int, int]{
				signal: states[3],
				time:   time + 2,
				label:  2,
			})
		},
		pinout: []string{"SET", "RESET", "INNER", "OUT"},
	},
}

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
