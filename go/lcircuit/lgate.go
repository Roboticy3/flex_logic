package lcircuit

import (
	c "flex-logic/collections"
	"flex-logic/ltypes"
)

/*
use a gate's state `[]S` at time `T` to schedule future events
*/
type Solver[S ltypes.LState, T ltypes.LTime] func([]S, T, *c.LEvents[S, T])

/*
define a type of gate that can be added to a circuit. Each type can have
multiple instances in a circuit.

`name` is the name of the gate type
`Solver` is the function that defines the gate's behavior
`pinout` is the list of pin names for the gate. Helps w/ implementing `Solver`
*/
type LGate[S ltypes.LState, T ltypes.LTime] struct {
	Name   string
	Solver Solver[S, T]
	Pinout []string
}

func (gtype LGate[S, T]) IsEmpty() bool {
	return gtype.Name == c.STRING_EMPTY
}
