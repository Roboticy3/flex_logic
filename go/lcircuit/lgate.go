package lcircuit

/*
 use a gate's state `[]S` at time `T` to schedule future events
*/
type solver[S lstate, T ltime] func([]S, T, *levents[S, T])

type lgate[S lstate, T ltime] struct {
	name   string
	solver solver[S, T]
}
