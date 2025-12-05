package lcircuit_test

import (
	c "flex-logic/collections"
	"flex-logic/lcircuit"
)

var testGates []lcircuit.LGate[int, int] = []lcircuit.LGate[int, int]{
	{
		Name: "AND",
		Solver: func(states []int, Time int, events *c.LEvents[int, int]) {
			states[2] = states[0] & states[1]
			events.Push(c.LEvent[int, int]{
				Signal: states[2],
				Time:   Time + 1,
				Label:  0,
			})
		},
		Pinout: []string{"A", "B", "OUT"},
	},
	{
		Name: "NOT",
		Solver: func(states []int, Time int, events *c.LEvents[int, int]) {
			states[1] = ^states[0]
			events.Push(c.LEvent[int, int]{
				Signal: states[1],
				Time:   Time + 1,
				Label:  1,
			})
		},
		Pinout: []string{"A", "OUT"},
	},
	{
		Name: "LATCH",
		Solver: func(states []int, Time int, events *c.LEvents[int, int]) {
			states[2] |= states[0]
			states[2] &^= states[1]
			states[3] = states[2]
			events.Push(c.LEvent[int, int]{
				Signal: states[3],
				Time:   Time + 2,
				Label:  2,
			})
		},
		Pinout: []string{"SET", "RESET", "INNER", "OUT"},
	},
}
