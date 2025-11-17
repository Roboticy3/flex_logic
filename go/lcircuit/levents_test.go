package lcircuit

import (
	"container/heap"
	"testing"
)

type testState string
type testTime int

func TestLeventsWithIntTime(t *testing.T) {
	events := &LEvents[testState, testTime]{}
	heap.Init(events)

	heap.Push(events, LEvent[testState, testTime]{time: 5, signal: "A", label: 1})
	heap.Push(events, LEvent[testState, testTime]{time: 3, signal: "B", label: 2})
	heap.Push(events, LEvent[testState, testTime]{time: 8, signal: "C", label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(LEvent[testState, testTime])
	if first.time != 3 {
		t.Errorf("Expected first event time to be 3, got %d", first.time)
	}

	second := heap.Pop(events).(LEvent[testState, testTime])
	if second.time != 5 {
		t.Errorf("Expected second event time to be 5, got %d", second.time)
	}

	third := heap.Pop(events).(LEvent[testState, testTime])
	if third.time != 8 {
		t.Errorf("Expected third event time to be 8, got %d", third.time)
	}
}

func TestLeventsWithFloatTime(t *testing.T) {
	events := &LEvents[testState, float64]{}
	heap.Init(events)

	heap.Push(events, LEvent[testState, float64]{time: 5.5, signal: "A", label: 1})
	heap.Push(events, LEvent[testState, float64]{time: 3.3, signal: "B", label: 2})
	heap.Push(events, LEvent[testState, float64]{time: 8.8, signal: "C", label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(LEvent[testState, float64])
	if first.time != 3.3 {
		t.Errorf("Expected first event time to be 3.3, got %f", first.time)
	}

	second := heap.Pop(events).(LEvent[testState, float64])
	if second.time != 5.5 {
		t.Errorf("Expected second event time to be 5.5, got %f", second.time)
	}

	third := heap.Pop(events).(LEvent[testState, float64])
	if third.time != 8.8 {
		t.Errorf("Expected third event time to be 8.8, got %f", third.time)
	}
}

func TestLeventsWithStringState(t *testing.T) {
	events := &LEvents[string, testTime]{}
	heap.Init(events)

	heap.Push(events, LEvent[string, testTime]{time: 10, signal: "ON", label: 1})
	heap.Push(events, LEvent[string, testTime]{time: 2, signal: "OFF", label: 2})
	heap.Push(events, LEvent[string, testTime]{time: 7, signal: "IDLE", label: 3})

	if events.Len() != 3 {
		t.Errorf("Expected 3 events, got %d", events.Len())
	}

	first := heap.Pop(events).(LEvent[string, testTime])
	if first.time != 2 {
		t.Errorf("Expected first event time to be 2, got %d", first.time)
	}

	second := heap.Pop(events).(LEvent[string, testTime])
	if second.time != 7 {
		t.Errorf("Expected second event time to be 7, got %d", second.time)
	}

	third := heap.Pop(events).(LEvent[string, testTime])
	if third.time != 10 {
		t.Errorf("Expected third event time to be 10, got %d", third.time)
	}
}
