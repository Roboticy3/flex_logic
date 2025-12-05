package lcircuit

import (
	c "flex-logic/collections"
	"flex-logic/ltypes"
)

/*
a pin on a specific component in a circuit

`valid` will be false by default, making new pins (eg `var zero LPin[S,T]`)
is an easy way to add empty entries.

`nets` will remain sorted and unique for easy search.
*/
type LPin[S ltypes.LState, T ltypes.LTime] struct {
	Nets  []c.Label
	Valid bool
}

func (pin LPin[S, T]) IsEmpty() bool {
	return !pin.Valid
}

func (pin LPin[S, T]) Len() int {
	return len(pin.Nets)
}

func (pin LPin[S, T]) Less(i, j int) bool {
	return pin.Nets[i] < pin.Nets[j]
}

func (pin LPin[S, T]) Swap(i, j int) {
	pin.Nets[i], pin.Nets[j] = pin.Nets[j], pin.Nets[i]
}

/*
A net between circuit pins. Drives connected nets with its state and processes
events with its `tid`.

`pins` will remain sorted and unique for easy search.
*/
type LNet[S ltypes.LState, T ltypes.LTime] struct {
	Pins  []c.Label
	Tid   c.Label
	State S
}

func (net LNet[S, T]) IsEmpty() bool {
	return len(net.Pins) == 0
}

func (net LNet[S, T]) Len() int {
	return len(net.Pins)
}

func (net LNet[S, T]) Less(i, j int) bool {
	return net.Pins[i] < net.Pins[j]
}

func (net LNet[S, T]) Swap(i, j int) {
	net.Pins[i], net.Pins[j] = net.Pins[j], net.Pins[i]
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
type LCircuit[S ltypes.LState, T ltypes.LTime] struct {
	netlist   *c.LLabeling[LNet[S, T]]
	pinlist   *c.LLabeling[LPin[S, T]]
	gatetypes *c.LLabeling[LGate[S, T]]
}

func (lc LCircuit[S, T]) GetNetlist() c.LLabeling[LNet[S, T]] {
	if lc.netlist == nil {
		return c.LLabeling[LNet[S, T]]{}
	}
	result := make(c.LLabeling[LNet[S, T]], len(*lc.netlist))
	for i, net := range *lc.netlist {
		if net.Pins != nil {
			result[i].Pins = make([]c.Label, len(net.Pins))
			copy(result[i].Pins, net.Pins)
		}
		result[i].Tid = net.Tid
		result[i].State = net.State
	}
	return result
}

func (lc LCircuit[S, T]) GetPinlist() c.LLabeling[LPin[S, T]] {
	if lc.pinlist == nil {
		return c.LLabeling[LPin[S, T]]{}
	}
	result := make(c.LLabeling[LPin[S, T]], len(*lc.pinlist))
	for i, pin := range *lc.pinlist {
		if pin.Nets != nil {
			result[i].Nets = make([]c.Label, len(pin.Nets))
			copy(result[i].Nets, pin.Nets)
		}
		result[i].Valid = pin.Valid
	}
	return result
}

func (lc LCircuit[S, T]) FindTypeName(gname string) c.Label {
	tid := c.LABEL_EMPTY
	for i := range *lc.gatetypes {
		if (*lc.gatetypes)[i].Name == gname {
			tid = c.Label(i)
			break
		}
	}

	if tid == c.LABEL_EMPTY {
		return c.LABEL_EMPTY
	}

	return tid
}

func (lc LCircuit[S, T]) GetGateTypes() c.LLabeling[LGate[S, T]] {
	if lc.gatetypes == nil {
		return c.LLabeling[LGate[S, T]]{}
	}
	result := make(c.LLabeling[LGate[S, T]], len(*lc.gatetypes))
	copy(result, *lc.gatetypes)
	return result
}

func (lc *LCircuit[S, T]) SetGateTypes(new_types c.LLabeling[LGate[S, T]]) {
	*lc.gatetypes = make(c.LLabeling[LGate[S, T]], len(new_types))
	copy(*lc.gatetypes, new_types)
}

func CreateCircuit[S ltypes.LState, T ltypes.LTime]() *LCircuit[S, T] {
	return &LCircuit[S, T]{
		&c.LLabeling[LNet[S, T]]{},
		&c.LLabeling[LPin[S, T]]{},
		&c.LLabeling[LGate[S, T]]{},
	}
}
