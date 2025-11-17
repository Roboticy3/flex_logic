package lcircuit

/*
a pin on a specific component in a circuit
*/
type LPin struct {
	component Label
	pin       Label
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
	return comp.tid == -1
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
type Lcircuit[S LState, T LTime] struct {
	nets_to_pins  LLabeling[LNet[S]]
	gates_to_nets LLabeling[LComponent]
	gtypes        []LGate[S, T]
	free_pins     LLabeling[Label]
}

/*
View and edit a circuit via gates.
*/
type LCGateController[S LState, T LTime] struct {
	*Lcircuit[S, T]
}

/*
 Add a gate of type `gtype` to the circuit with no connections, returning the
 gate id

 O(n + m) for n gates and m nets.
*/
func (gview LCGateController[S, T]) AddGate(gname string) Label {
	//With no connections, each pin in a gate will induce a separate net on the
	//	circuit.

	tid := Label(-1)
	for i, gt := range gview.gtypes {
		if gt.name == gname {
			tid = Label(i)
			break
		}
	}

	if tid == -1 {
		return -1 // Gate type not found
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
				pin:       Label(i),
			}},
			state: zero,
		}, int(nid))
		gview.gates_to_nets[gid].nets[i] = nid
	}

	return gid
}

/*
 Remove a gate with id `gid` from the circuit. Returns the found gate type if
 successful, or -1

 O(p) for p pins on the gate.
*/
func (gview LCGateController[S, T]) RemoveGate(gid Label) Label {

	if gview.gates_to_nets.Get(gid) == nil {
		return -1
	}

	empty_net := LNet[S]{}
	empty_component := LComponent{-1, []Label{}}

	for _, nid := range gview.gates_to_nets[gid].nets {
		gview.nets_to_pins.Remove(nid, empty_net)
	}

	result := gview.gates_to_nets[gid].tid
	gview.gates_to_nets.Remove(gid, empty_component)

	return result
}

/*
 Get the type of the gate with id `gid`

 If no gate is found, a gate type passing IsEmpty is returned.
 Breaks if gid is out of range.

 O(1)
*/
func (gview LCGateController[S, T]) GetGateName(gid Label) string {
	return gview.gtypes[gview.gates_to_nets[gid].tid].name
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
	*Lcircuit[S, T]
}

func (pview LCPinController[S, T]) AddPin() {

	//Add a blank net pointing past the end of the array
	var zero S
	nid := pview.nets_to_pins.Add(LNet[S]{
		pins: []LPin{{
			component: -1,
			pin:       Label(len(pview.free_pins)),
		}},
		state: zero,
	}, 0)

	pview.free_pins.Add(nid, 0)
}
