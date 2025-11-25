package lcircuit

import (
	"slices"
	"sort"
)

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

O(uq) for q, u average connections on nets, pins. In a standard digital logic
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

O(q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) Detach(nid Label, pid Label) bool {

	p_pin, p_net := nc.pinlist.Get(pid), nc.netlist.Get(nid)
	if p_pin == nil || p_net == nil {
		return false
	}

	//remove nid from sorted array. Empty arrays are considered empty entries
	i := sort.Search(len(p_pin.nets), func(k int) bool { return p_pin.nets[k] >= nid })
	p_pin.nets = append(p_pin.nets[0:i], p_pin.nets[i+1:]...)

	//remove pid from sorted array
	j := sort.Search(len(p_net.pins), func(k int) bool { return p_net.pins[k] >= pid })
	p_net.pins = append(p_net.pins[0:j], p_net.pins[j+1:]...)

	return true
}

/*
Attach net `nid` to pin `pid`. If either do not exist, or the connection
already exists, return false.

O(q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) Attach(nid Label, pid Label) bool {

	p_pin, p_net := nc.pinlist.Get(pid), nc.netlist.Get(nid)
	if p_pin == nil || p_net == nil {
		return false
	}

	//insert nid into sorted array
	i := sort.Search(len(p_pin.nets), func(k int) bool { return p_pin.nets[k] >= nid })
	p_pin.nets = slices.Insert(p_pin.nets, i, nid)

	//insert pid into sorted array
	j := sort.Search(len(p_net.pins), func(k int) bool { return p_net.pins[k] >= pid })
	p_net.pins = slices.Insert(p_net.pins, j, pid)

	return true
}

/*
Merge nets at `nids` together. The first valid element of `nids` is the target
net for the merge, so its `state` and `tid` will dominate.

Return the merge target, or `LABEL_EMPTY` if no target was found

O(nq log q) for n valid elements of `nids`, q average connections on a net.
*/
func (nc LCNetController[S, T]) Merge(nids []Label) Label {

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
		return LABEL_EMPTY
	}

	if target_i == len(nids)-1 {
		return LABEL_EMPTY
	}

	nids = nids[target_i+1:]

	//merge each valid net after target into target
	for _, nid := range nids {
		nc.MergeTwo(target, nid)
	}

	return target
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
		//Since net pins are sorted, the merge is O(q) with this:
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

/*
Add a net to the circuit. Copies the input into the circuit, sorting and
validating pins.

If no valid pins are on net.pins, the function fails and returns LABEL_EMPTY

O(q^2) for q average connections
*/
func (nc LCNetController[S, T]) Add(net LNet[S, T]) Label {

	//new_net to copy net into. Helps with checking pins
	new_net := LNet[S, T]{
		pins:  []Label{},
		tid:   net.tid,
		state: net.state,
	}
	nid := nc.netlist.Add(new_net, 0)

	//attempt to attach the copy to all requested pins.
	//if no valid attachments are found, valid remains false
	valid := false
	for _, pid := range net.pins {
		if nc.Attach(nid, pid) {
			valid = true
		}
	}

	//if no valid pins were found, undo the addition and return LABEL_EMPTY
	//	to indicate failure.
	if !valid {
		var zero LNet[S, T]
		nc.netlist.Remove(nid, zero)
		return LABEL_EMPTY
	}

	return nid
}

/*
Remove a net from the circuit. Disconnects all pins in net.pins from each other.

net.pins is searched in order. If any two pins are found to be connected by a
net `nid`, `nid` is detached from the later pin.

If `empty_only` is true, only nets of type `LABEL_EMPTY` will be detached.

# Returns true if nets were disconnected from each other, and false otherwise

O(nq^2) for n pins with q average connections. Only happens when the set of pins
is maximally connected/complete.
*/
func (nc LCNetController[S, T]) Remove(net LNet[S, T], empty_only bool) bool {
	//build neighborhoods as sets for faster lookup
	neighborhood := map[Label]bool{}
	result := false

	for _, pid := range net.pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		//enforce that the nid is detached.
		for _, nid := range p_pin.nets {
			p_net := nc.netlist.Get(nid)
			if p_net == nil {
				continue
			}

			_, ok := neighborhood[nid]
			if ok {
				//split pid from nid
				nc.Detach(nid, pid)
				result = true
			} else if (!empty_only) || p_net.tid == LABEL_EMPTY {
				//otherwise, track the nid for later
				neighborhood[nid] = true
			}
		}
	}

	return result
}

/*
Wire 2 valid tid'd nets with a new net of tid `LABEL_EMPTY`. If a net of type
`LABEL_EMPTY` already exists on either pin, merge all such nets together.

Returns the label of the created net, or `LABEL_EMPTY` if either pin was invalid

O(q^2 log q) for q average connections
*/
func (nc LCNetController[S, T]) AddWire(pid1 Label, pid2 Label) Label {

	//check for nets to merge into on the pins, validating that the pins exist
	//in the process
	p_pins := [](*LPin[S, T]){nc.pinlist.Get(pid1), nc.pinlist.Get(pid2)}
	if p_pins[0] == nil || p_pins[1] == nil {
		return LABEL_EMPTY
	}

	//O(q)
	merges := []Label{}
	for _, p_pin := range p_pins {
		for _, nid := range p_pin.nets {
			p_net := nc.netlist.Get(nid)
			if p_net == nil {
				continue
			}

			if p_net.tid != -1 {
				merges = append(merges, nid)
			}
		}
	}

	//Add the new net
	var zero S
	wire := LNet[S, T]{[]Label{pid1, pid2}, LABEL_EMPTY, zero}
	//O(q^2)
	nid := nc.Add(wire)

	if nid == LABEL_EMPTY {
		return LABEL_EMPTY //should never happen if pin check passed
	}

	//If merges were found, merge the discovered nets together and return the
	//merge target
	//O(nq log q) and n is approximately q
	if len(merges) > 0 {
		merges = append(merges, nid)

		return nc.Merge(merges)
	}

	//Otherwise, just return the id of the added net.
	return nid
}

/*
Disconnect `pid1` from `pid2` along nets of type `LABEL_EMPTY` (won't
disconnect two pins from a gate).

O(q^2)
*/
func (nc LCNetController[S, T]) RemoveWire(pid1 Label, pid2 Label) bool {
	var wire LNet[S, T]
	wire.pins = []Label{pid1, pid2}
	return nc.Remove(wire, true)
}
