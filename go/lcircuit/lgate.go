package lcircuit

import (
	c "flex-logic/collections"
	"flex-logic/ltypes"
)

/*
use a gate's state from pins P to add to event queue LEvents[S, T]
*/
type Solver[P, S ltypes.LState, T ltypes.LTime] func(P, T, *c.LEvents[S, T])

/*
define a type of gate that can be added to a circuit. Each type can have
multiple instances in a circuit.

`name` is the name of the gate type
`Solver` is the function that defines the gate's behavior
`pinout` is the list of pin names for the gate. Helps w/ implementing `Solver`
*/
type LGate[P, S ltypes.LState, T ltypes.LTime] struct {
	Name   string
	Solver Solver[P, S, T]
	Pinout []string
}

func (gtype LGate[P, S, T]) IsEmpty() bool {
	return gtype.Name == c.STRING_EMPTY
}
