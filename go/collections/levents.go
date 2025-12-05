package collections

import (
	"container/heap"
	"flex-logic/ltypes"
)

/*
Event specifying a `signal` at a pin `label` occuring at time `time`
*/
type LEvent[S ltypes.LState, T ltypes.LTime] struct {
	Time   T
	Signal S
	Label  Label
}

type LEvents[S ltypes.LState, T ltypes.LTime] []LEvent[S, T]

var _ heap.Interface = (*LEvents[any, int])(nil)

// Implement heap interface
//
//	See reference implemention at https://pkg.go.dev/container/heap
func (events LEvents[S, T]) Len() int { return len(events) }
func (events LEvents[S, T]) Less(i, j int) bool {
	return events[i].Time < events[j].Time
}
func (events LEvents[S, T]) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}
func (events *LEvents[S, T]) Push(x any) {
	*events = append(*events, x.(LEvent[S, T]))
}
func (events *LEvents[S, T]) Pop() any {
	old := *events
	n := len(old)
	x := old[n-1]
	*events = old[0 : n-1]
	return x
}
