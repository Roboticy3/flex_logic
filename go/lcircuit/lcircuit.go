package lcircuit

/*
a pin on a specific component in a circuit
*/
type LPin struct {
	component Label
	pin       int
}

type LNet[S LState] struct {
	pins  []LPin
	state S
}

func (net LNet[S]) IsEmpty() bool {
	return len(net.pins) == 0
}

/*
a component in the circuit has a gate and connections to nets

nets is ordered to match `gtype.pinout`
*/
type LComponent struct {
	tid  Label
	nets []Label
}

func (comp LComponent) IsEmpty() bool {
	return comp.tid == LABEL_EMPTY
}

/*
 Base type for a circuit
	`nets_to_pins`:		mapping of net states to arrays of pins
	`gates_to_nets`:	mapping of gates to nets that drive them/that they drive
	`gtypes`:					collection of immutable gate types that are allowed for
		dynamic changes
	`free_pins`:			mutable pinout of the current circuit associated with net
		ids
*/
type LCircuit[S LState, T LTime] struct {
	nets_to_pins  LLabeling[LNet[S]]
	gates_to_nets LLabeling[LComponent]
	gtypes        []LGate[S, T]
	free_pins     LLabeling[Label]
}

/*
View and edit a circuit via gates.
*/
type LCGateController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
 Add a gate of type `gtype` to the circuit with no connections, returning the
 gate id

 O(n + m) for n gates and m nets.
*/
func (gview LCGateController[S, T]) AddGate(gname string) Label {
	//With no connections, each pin in a gate will induce a separate net on the
	//	circuit.

	tid := Label(LABEL_EMPTY)
	for i, gt := range gview.gtypes {
		if gt.name == gname {
			tid = Label(i)
			break
		}
	}

	if tid == LABEL_EMPTY {
		return LABEL_EMPTY // Gate type not found
	}

	// Add the gate component
	gid := gview.gates_to_nets.Add(LComponent{
		tid:  tid,
		nets: make([]Label, len(gview.gtypes[tid].pinout)),
	}, 0)

	// Sweep through nets to add, binding each one to their respective pin
	nid := Label(0)
	var zero S
	for i := range gview.gtypes[tid].pinout {
		nid = gview.nets_to_pins.Add(LNet[S]{
			pins: []LPin{{
				component: gid,
				pin:       i,
			}},
			state: zero,
		}, int(nid))
		gview.gates_to_nets[gid].nets[i] = nid
	}

	return gid
}

/*
 Remove a gate with id `gid` from the circuit. Returns the found gate type if
 successful, or LABEL_EMPTY

 O(pq) for p pins on the gate with q expected connections
*/
func (gview LCGateController[S, T]) RemoveGate(gid Label) Label {

	component := gview.gates_to_nets.Get(gid)
	if component == nil {
		return LABEL_EMPTY
	}

	//detach each pin of the gate from its network
	for i := range component.nets {
		LCNetView[S, T](gview).Detach(LPin{gid, i})
	}

	//then remove it
	empty_component := LComponent{LABEL_EMPTY, []Label{}}
	gview.gates_to_nets.Remove(gid, empty_component)

	result := component.tid
	return result
}

/*
 Get the type of the gate with id `gid`

 If no gate is found, a gate type passing IsEmpty is returned.
 Breaks if gid is out of range.

 O(1)
*/
func (gview LCGateController[S, T]) GetGateName(gid Label) string {
	component := gview.gates_to_nets.Get(gid)
	if component == nil {
		return "?Missing?"
	}
	return gview.gtypes[component.tid].name
}

/*
 List gates currently in the circuit by `gid`

 O(n) for n gates.
*/
func (gview LCGateController[S, T]) ListGates() []Label {
	var gids []Label
	for i := 0; i < gview.gates_to_nets.Len(); i++ {
		if !(gview.gates_to_nets[i].IsEmpty()) {
			gids = append(gids, Label(i))
		}
	}
	return gids
}

type LCPinController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
 Add 1 pin to the circuit and induce a net

 O(1)
*/
func (pview LCPinController[S, T]) AddFreePin() {

	//Add a blank net pointing past the end of the array
	var zero S
	nid := pview.nets_to_pins.Add(LNet[S]{
		pins: []LPin{{
			component: LABEL_EMPTY,
			pin:       len(pview.free_pins),
		}},
		state: zero,
	}, 0)

	//Fill in the pin being referenced by the net we just made
	pview.free_pins.Add(nid, 0)
}

/*
 Remove 1 pin and disconnect it from its net

 O(q) for q expected connections from `pid`
*/
func (pview LCPinController[S, T]) RemoveFreePin(pid Label) {

	//detach the pin from its network
	LCNetView[S, T](pview).Detach(LPin{
		component: LABEL_EMPTY,
		pin:       int(pid),
	})

	//then remove it
	pview.free_pins.Remove(pid, LABEL_EMPTY)
}

/*
 Add and remove connections to existing nets
 Since nets should always have at least 1 pin, and adding pins induces nets onto
	the circuit, adding nets is not meaningful. Only merging sets of pins or
	splitting a net label
*/
type LCNetView[S LState, T LTime] struct {
	*LCircuit[S, T]
}

func (nview LCNetView[S, T]) Detach(pin LPin) {
	nid := nview.gates_to_nets[pin.component].nets[pin.pin]

	if nid.IsEmpty() {
		return
	}

	//swap out the pin to the end and shrink the slice header
	pins := nview.nets_to_pins[nid].pins
	for i, p := range pins {
		if p == pin {
			pins[i] = pins[len(pins)-1]
			pins = pins[:len(pins)-1]
		}
	}

	//if the netlist is emptied, clear it from the labeling.
	var zero LNet[S]
	if len(pins) == 0 {
		nview.nets_to_pins.Remove(nid, zero)
	} else {
		//otherwise copy in the updated slice header
		nview.nets_to_pins[nid].pins = pins
	}

}
