package lcircuit

import (
	"sort"
)

/*
a pin on a specific component in a circuit

`valid` will be false by default, making new pins (eg `var zero LPin[S,T]`)
is an easy way to add empty entries.

`nets` will remain sorted for easy search.
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

`pins` will remain sorted for easy merging.
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

type LCPinController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
Add a pin to the circuit.

If the provided `nid` is `LABEL_EMPTY`, the pin is added with no associated
net.

If `nid` is a valid label, but doesn't correspond to a valid net, the function
fails and returns `LABEL_EMPTY`

If `nid` corresponds to a valid net, the pin is added to the circuit and
connected to the net.

O(q log q) for q average connections on a net.
*/
func (pc LCPinController[S, T]) AddPin(nid Label) Label {

	if nid == LABEL_EMPTY {
		return pc.pinlist.Add(LPin[S, T]{[]Label{}, true}, 0)
	}

	p_net := pc.netlist.Get(nid)
	if p_net == nil {
		return LABEL_EMPTY
	}

	result := pc.pinlist.Add(LPin[S, T]{[]Label{nid}, true}, 0)
	for _, pid := range p_net.pins {
		if pid == result {
			return result
		}
	}

	p_net.pins = append(p_net.pins, result)
	sort.Sort(p_net)

	return result
}

/*
Remove a pin at `pid` from the circuit. Returns true if the pin exists.
False otherwise.

The pin is disconnected from any associated nets. If these nets become fully
disconnected, they are removed.

O(uq + u log q) for q, u average connections on nets, pins. In a standard digital logic
circuit, pins should not exceed two nets- one component and one wire.
*/
func (pc LCPinController[S, T]) RemovePin(pid Label) bool {

	p_pin := pc.pinlist.Get(pid)
	if p_pin == nil {
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

func (pc LCPinController[S, T]) ListPins() []Label {
	result := []Label{}
	for pid, pin := range pc.pinlist {
		if !pin.IsEmpty() {
			result = append(result, Label(pid))
		}
	}
	return result
}

type LCNetController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
Detach net `nid` from pin `pid`

If detaching this pin would leave net `nid` with zero connections, the net
is removed from the circuit.

O(q + log q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) Detach(nid Label, pid Label) bool {

	p_net := nc.netlist.Get(nid)

	if len(p_net.pins) == 1 && p_net.pins[0] == pid {
		var zero LNet[S, T]
		nc.netlist.Remove(nid, zero)
		return true
	}

	//binary search O(log q)
	i := sort.Search(len(p_net.pins), func(i int) bool { return p_net.pins[i] >= pid })
	if i < 0 || i >= len(p_net.pins) {
		return false
	}

	//lazy remove (O(q))
	p_net.pins = append(p_net.pins[:i], p_net.pins[i+1:]...)

	return true
}

/*
Attach net `nid` to pin `pid`. If either do not exist, or the connection
already exists, return false.

O(log q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) Attach(nid Label, pid Label) bool {
	return false
}

/*
Merge nets at `nids` together. The first valid element of `nids` is the target
net for the merge, so its `state` and `tid` will dominate.

O(nq log q) for n valid elements of `nids`, q average connections on a net.
*/
func (nc LCNetController[S, T]) Merge(nids []Label) {

	target := LABEL_EMPTY
	target_i := -1 //store the index where the target was found for slicing
	for i, nid := range nids {
		p_net := nc.netlist.Get(nid)
		if p_net != nil {
			target = nid
			target_i = i
			break
		}
	}

	if target == LABEL_EMPTY {
		return
	}

	if target_i == len(nids)-1 {
		return
	}

	nids = nids[target_i+1:]

	//merge each valid net after target into target
	for _, nid := range nids {
		nc.MergeTwo(target, nid)
	}
}

/*
Merge net nid2 into nid1. If nid2 is invalid, return false. If nid1 is
invalid, replace it with nid2.

O(q log q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) MergeTwo(nid1 Label, nid2 Label) bool {
	p_source := nc.netlist.Get(nid2)
	if p_source == nil {
		return false
	}

	p_target := nc.netlist.Get(nid1)
	if p_target == nil {
		if nid1 == LABEL_EMPTY {
			return false
		}

		//In this case, just make a new target with no connections.
		new_target := LNet[S, T]{
			p_source.pins,
			p_source.tid,
			p_source.state,
		}
		nc.netlist[nid1] = new_target
	} else {
		//Otherwise, merge the pin arrays, use target for the tid and state
		//Since net pins are sorted, the merge is linear with interlacing:
		pins := make([]Label, len(p_target.pins)+len(p_source.pins))
		i, j := 0, 0
		for p := range pins {
			//push target
			if i < len(p_target.pins) && p_target.pins[i] < p_source.pins[j] {
				pins[p] = p_target.pins[i]
				i++
				//remove duplicates
			} else if i < len(p_target.pins) && j < len(p_source.pins) && p_target.pins[i] == p_source.pins[j] {
				pins[p] = p_target.pins[i]
				pins = pins[:len(pins)-1]
				i++
				j++
				//push source
			} else if j < len(p_source.pins) {
				pins[p] = p_source.pins[j]
				j++
			}
		}

		//set pins, but keep the target's tid and state
		p_target.pins = pins
	}

	//retarget all pins pointing at source to point at target O(q log q)
	for _, pid := range p_target.pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		//binary search O(log q)
		i := sort.Search(len(p_pin.nets), func(i int) bool { return p_pin.nets[i] >= nid2 })
		if i < 0 || i >= len(p_pin.nets) {
			continue
		}

		//retarget the pin to point at target
		p_pin.nets[i] = nid1
	}

	//remove the source *without* removing pins, since they have been retargeted
	nc.netlist.Remove(nid2, LNet[S, T]{})

	return true
}
