package lcircuit

/*
View and edit a circuit via gates.
*/
type LCGateController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
Adding a gate of type `gname` induces pins and nets onto the circuit.

O(t + p) for t types and p pins.
*/
func (gc LCGateController[S, T]) AddGate(gname string) Label {

	//Find the type label of the gate
	tid := LCGateTypeController[S, T](gc).FindTypeName(gname)
	if tid == LABEL_EMPTY {
		return LABEL_EMPTY
	}

	//Add in the net first, since its easier to fill in the pins on one net
	//	than to fill in the net on many pins
	//Might change memory management later to take slices out of a flat array here
	pincount := len(gc.gatetypes[tid].pinout)
	pins := make([]Label, pincount)
	var zero S
	nid := gc.netlist.Add(LNet[S, T]{
		pins,
		tid,
		zero,
	}, 0)

	//Drizzle new pins into the pinlist. The new net's pin array will be sorted
	//as a result :)
	pid := 0
	for i := range pincount {
		pins[i] = gc.pinlist.Add(LPin[S, T]{[]Label{nid}, true}, pid)
		pid = int(pins[i]) + 1
	}

	return nid
}

/*
Remove the gate at `gid`. Returns true if the gate existed and false otherwise

O(pqu) for p pins, q, u average connections on nets, pins.
*/
func (gc LCGateController[S, T]) RemoveGate(gid Label) bool {

	//Find a gate. If the net id belongs to a wire cluster or is empty, ignore it
	p_net := gc.netlist.Get(gid)
	if p_net.tid == -1 || p_net.IsEmpty() {
		return false
	}

	//Otherwise, we want to remove all associated pins, then the net itself from
	//	their respective labelings.
	//Removing a pin may also have effects on other nets. Use LCPinController
	//	to handle this.
	//Removing the last pin will fully disconnect this net and remove it auto-
	//	matically, so we shouldn't have to worry about it  after this.
	for len(p_net.pins) > 0 {
		pid := p_net.pins[0]
		LCPinController[S, T](gc).RemovePin(pid)
	}

	return true
}

/*
List valid gate ids

O(n) for n gates.
*/
func (gc LCGateController[S, T]) ListGateIds() []Label {
	result := []Label{}
	for nid, net := range gc.netlist {
		if !net.IsEmpty() && net.tid != LABEL_EMPTY {
			result = append(result, Label(nid))
		}
	}

	return result
}
