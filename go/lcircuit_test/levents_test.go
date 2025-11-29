package lcircuit_test

import (
	"container/heap"
	"flex-logic/lcircuit"
	"testing"
)

type testState string
type testTime int

func TestLeventsWithIntTime(t *testing.T) {
	events := &lcircuit.LEvents[testState, testTime]{}
	heap.Init(events)

	heap.Push(events, lcircuit.LEvent[testState, testTime]{Time: 5, Signal: "A", Label: 1})
	heap.Push(events, lcircuit.LEvent[testState, testTime]{Time: 3, Signal: "B", Label: 2})
	heap.Push(events, lcircuit.LEvent[testState, testTime]{Time: 8, Signal: "C", Label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(lcircuit.LEvent[testState, testTime])
	if first.Time != 3 {
		t.Errorf("Expected first event Time to be 3, got %d", first.Time)
	}

	second := heap.Pop(events).(lcircuit.LEvent[testState, testTime])
	if second.Time != 5 {
		t.Errorf("Expected second event Time to be 5, got %d", second.Time)
	}

	third := heap.Pop(events).(lcircuit.LEvent[testState, testTime])
	if third.Time != 8 {
		t.Errorf("Expected third event Time to be 8, got %d", third.Time)
	}
}

func TestLeventsWithFloatTime(t *testing.T) {
	events := &lcircuit.LEvents[testState, float64]{}
	heap.Init(events)

	heap.Push(events, lcircuit.LEvent[testState, float64]{Time: 5.5, Signal: "A", Label: 1})
	heap.Push(events, lcircuit.LEvent[testState, float64]{Time: 3.3, Signal: "B", Label: 2})
	heap.Push(events, lcircuit.LEvent[testState, float64]{Time: 8.8, Signal: "C", Label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(lcircuit.LEvent[testState, float64])
	if first.Time != 3.3 {
		t.Errorf("Expected first event Time to be 3.3, got %f", first.Time)
	}

	second := heap.Pop(events).(lcircuit.LEvent[testState, float64])
	if second.Time != 5.5 {
		t.Errorf("Expected second event Time to be 5.5, got %f", second.Time)
	}

	third := heap.Pop(events).(lcircuit.LEvent[testState, float64])
	if third.Time != 8.8 {
		t.Errorf("Expected third event Time to be 8.8, got %f", third.Time)
	}
}

func TestLeventsWithStringState(t *testing.T) {
	events := &lcircuit.LEvents[string, testTime]{}
	heap.Init(events)

	heap.Push(events, lcircuit.LEvent[string, testTime]{Time: 10, Signal: "ON", Label: 1})
	heap.Push(events, lcircuit.LEvent[string, testTime]{Time: 2, Signal: "OFF", Label: 2})
	heap.Push(events, lcircuit.LEvent[string, testTime]{Time: 7, Signal: "IDlcircuit.LE", Label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(lcircuit.LEvent[string, testTime])
	if first.Time != 2 {
		t.Errorf("Expected first event Time to be 2, got %d", first.Time)
	}

	second := heap.Pop(events).(lcircuit.LEvent[string, testTime])
	if second.Time != 7 {
		t.Errorf("Expected second event Time to be 7, got %d", second.Time)
	}

	third := heap.Pop(events).(lcircuit.LEvent[string, testTime])
	if third.Time != 10 {
		t.Errorf("Expected third event Time to be 10, got %d", third.Time)
	}
}
