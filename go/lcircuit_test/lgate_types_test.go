package lcircuit_test

import (
	c "flex-logic/collections"
	"flex-logic/lcircuit"
)

var testGates []lcircuit.LGate[[]lcircuit.LPin[int, int], int, int] = []lcircuit.LGate[[]lcircuit.LPin[int, int], int, int]{
	{
		Name: "AND",
		Solver: func(pins []lcircuit.LPin[int, int], Time int, events *c.LEvents[int, int]) {
			pins[2].State = pins[0].State & pins[1].State
			events.Push(c.LEvent[int, int]{
				Signal: pins[2].State,
				Time:   Time + 1,
				Label:  0,
			})
		},
		Pinout: []string{"A", "B", "OUT"},
	},
	{
		Name: "NOT",
		Solver: func(pins []lcircuit.LPin[int, int], Time int, events *c.LEvents[int, int]) {
			pins[1].State = ^pins[0].State
			events.Push(c.LEvent[int, int]{
				Signal: pins[1].State,
				Time:   Time + 1,
				Label:  1,
			})
		},
		Pinout: []string{"A", "OUT"},
	},
	{
		Name: "LATCH",
		Solver: func(pins []lcircuit.LPin[int, int], Time int, events *c.LEvents[int, int]) {
			pins[2].State |= pins[0].State
			pins[2].State &^= pins[1].State
			pins[3].State = pins[2].State
			events.Push(c.LEvent[int, int]{
				Signal: pins[3].State,
				Time:   Time + 2,
				Label:  2,
			})
		},
		Pinout: []string{"SET", "RESET", "INNER", "OUT"},
	},
}
