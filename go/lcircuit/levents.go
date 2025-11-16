package Lcircuit

import (
	"container/heap"
)

/*
Event specifying a `signal` at a pin `label` occuring at time `time`
*/
type levent[S Lstate, T Ltime] struct {
	time   T
	signal S
	label  int
}

type levents[S Lstate, T Ltime] []levent[S, T]

var _ heap.Interface = (*levents[any, int])(nil)

// Implement heap interface
//
//	See reference implemention at https://pkg.go.dev/container/heap
func (events levents[S, T]) Len() int { return len(events) }
func (events levents[S, T]) Less(i, j int) bool {
	return events[i].time < events[j].time
}
func (events levents[S, T]) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}
func (events *levents[S, T]) Push(x any) {
	*events = append(*events, x.(levent[S, T]))
}
func (events *levents[S, T]) Pop() any {
	old := *events
	n := len(old)
	x := old[n-1]
	*events = old[0 : n-1]
	return x
}
