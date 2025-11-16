package lcircuit

var testGates []Lgate[int, int] = []Lgate[int, int]{
	{
		name: "AND",
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
		name: "NOT",
		Solver: func(states []int, time int, events *levents[int, int]) {
			states[1] = ^states[0]
			events.Push(levent[int, int]{
				signal: states[1],
				time:   time + 1,
				label:  1,
			})
		},
		pinout: []string{"A", "OUT"},
	},
	{
		name: "LATCH",
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
