package lcircuit

/*
a pin on a specific component in a circuit

`valid` will be false by default, making new pins (eg `var zero LPin[S,T]`)
is an easy way to add empty entries.

`nets` will remain sorted and unique for easy search.
*/
type LPin[S LState, T LTime] struct {
	nets  []Label
	valid bool
}

func (pin LPin[S, T]) IsEmpty() bool {
	return !pin.valid
}

func (pin LPin[S, T]) Len() int {
	return len(pin.nets)
}

func (pin LPin[S, T]) Less(i, j int) bool {
	return pin.nets[i] < pin.nets[j]
}

func (pin LPin[S, T]) Swap(i, j int) {
	pin.nets[i], pin.nets[j] = pin.nets[j], pin.nets[i]
}

/*
A net between circuit pins. Drives connected nets with its state and processes
events with its `tid`.

`pins` will remain sorted and unique for easy search.
*/
type LNet[S LState, T LTime] struct {
	pins  []Label
	tid   Label
	state S
}

func (net LNet[S, T]) IsEmpty() bool {
	return len(net.pins) == 0
}

func (net LNet[S, T]) Len() int {
	return len(net.pins)
}

func (net LNet[S, T]) Less(i, j int) bool {
	return net.pins[i] < net.pins[j]
}

func (net LNet[S, T]) Swap(i, j int) {
	net.pins[i], net.pins[j] = net.pins[j], net.pins[i]
}

/*
	Base type for a circuit
		`netlist`	labeling over nets on the circuit
		`pinlist` labeling over pins on the circuit
		`gatetypes` labeling over the `tid` field of each net in `netlist`

		Each pin belongs to n nets
		Each net sees p pins and has an optional gate type.

	Edit using only LCGateController and LCWireController to maintain a valid
	state for simulation. Otherwise simulating a circuit can have undefined
	behavior.
*/
type LCircuit[S LState, T LTime] struct {
	netlist   LLabeling[LNet[S, T]]
	pinlist   LLabeling[LPin[S, T]]
	gatetypes LLabeling[LGate[S, T]]
}

type LCGateTypeController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

func (gtc LCGateTypeController[S, T]) FindTypeName(gname string) Label {
	tid := LABEL_EMPTY
	for i := range gtc.gatetypes {
		if gtc.gatetypes[i].name == gname {
			tid = Label(i)
			break
		}
	}

	if tid == LABEL_EMPTY {
		return LABEL_EMPTY
	}

	return tid
}
