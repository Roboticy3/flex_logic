package lcircuit

/*
	a pin on a specific component in a circuit

	`valid` will be false by default, making new pins (eg `var zero LPin[S,T]`)
	is an easy way to add empty entries.
*/
type LPin[S LState, T LTime] struct {
	nets  []Label
	valid bool
}

func (pin LPin[S, T]) IsEmpty() bool {
	return !pin.valid
}

type LNet[S LState, T LTime] struct {
	pins  []Label
	tid   Label
	state S
}

func (net LNet[S, T]) IsEmpty() bool {
	return len(net.pins) == 0
}

/*
 Base type for a circuit
	`netlist`	labeling over nets on the circuit
	`pinlist` labeling over pins on the circuit

	Each pin belongs to n nets
	Each net sees p pins and has an optional gate type.
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

	pid := 0
	for i := range pincount {
		pins[i] = gc.pinlist.Add(LPin[S, T]{[]Label{nid}, true}, pid)
	}

	return nid
}

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
	for _, pid := range p_net.pins {
		LCPinController[S, T](gc).RemovePin(pid)
	}

	return true
}

func (gc LCGateController[S, T]) ListGateIds() []Label {
	result := []Label{}
	for nid, net := range gc.netlist {
		if !net.IsEmpty() && net.tid != LABEL_EMPTY {
			result = append(result, Label(nid))
		}
	}

	return result
}

type LCPinController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

func (pc LCPinController[S, T]) RemovePin(pid Label) bool {

	p_pin := pc.pinlist.Get(pid)
	if p_pin.IsEmpty() {
		return false
	}

	//find this pin in the connected nets and remove it
	for _, nid := range p_pin.nets {
		LCNetController[S, T](pc).Detach(nid, pid)
	}

	var zero LPin[S, T]
	pc.pinlist.Remove(pid, zero)

	return true
}

type LCNetController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
	Detach net `nid` from pin `pid`

	If detaching this pin would leave net `nid` with zero connections, the net
		is removed from the circuit.
*/
func (nc LCNetController[S, T]) Detach(nid Label, pid Label) bool {

	p_net := nc.netlist.Get(nid)

	if len(p_net.pins) == 1 && p_net.pins[0] == pid {
		var zero LNet[S, T]
		nc.netlist.Remove(nid, zero)
		return true
	}

	for i, pid2 := range p_net.pins {
		if pid == pid2 {
			//lazy remove (O(1))
			p_net.pins[i] = p_net.pins[len(p_net.pins)-1]
			p_net.pins = p_net.pins[:len(p_net.pins)-1]
			return true
		}
	}

	return false
}
