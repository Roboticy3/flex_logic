package lcircuit

import (
	"slices"
	"sort"
)

/*
View and edit a circuit via its nets, also with more advanced primitives for
wires.
*/
type LCNetController[S LState, T LTime] struct {
	*LCircuit[S, T]
}

/*
Get pins attached to `nid`, or empty array if `nid` is invalid

O(q)
*/
func (nc LCNetController[S, T]) GetPins(nid Label) []Label {
	p_net := nc.netlist.Get(nid)
	if p_net == nil {
		return []Label{}
	}

	pins := make([]Label, len(p_net.Pins))
	copy(pins, p_net.Pins)
	return pins
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
	i := sort.Search(len(p_pin.Nets), func(k int) bool { return p_pin.Nets[k] >= nid })
	p_pin.Nets = append(p_pin.Nets[0:i], p_pin.Nets[i+1:]...)

	//remove pid from sorted array
	j := sort.Search(len(p_net.Pins), func(k int) bool { return p_net.Pins[k] >= pid })
	p_net.Pins = append(p_net.Pins[0:j], p_net.Pins[j+1:]...)

	//if the net is emptied out, don't have to explicitly remove since a net with
	//no connections is invalid implicitly.

	return true
}

/*
Internal function. Produce one-way connection from `nid` to `pid`.

Assumes `p_net` is not `nil` and corresponds to `nid`. Returns false if `nid`
was already connected to `pid`.
*/
func (nc LCNetController[S, T]) attach_net(nid Label, pid Label, p_net *LNet[S, T]) bool {
	i := sort.Search(len(p_net.Pins), func(k int) bool { return p_net.Pins[k] >= pid })
	if i >= len(p_net.Pins) || p_net.Pins[i] != pid {
		p_net.Pins = slices.Insert(p_net.Pins, i, pid)
		return true
	}

	return false
}

/*
Internal function. Produce one-way connection from `pid` to `nid`.

Assumes `p_pin` is not `nil` and corresponds to `pid`. Returns false if `pid`
was already connected to `nid`.
*/
func (nc LCNetController[S, T]) attach_pin(nid Label, pid Label, p_pin *LPin[S, T]) bool {
	i := sort.Search(len(p_pin.Nets), func(k int) bool { return p_pin.Nets[k] >= nid })
	//skip insertion if the pid already exists
	if i >= len(p_pin.Nets) || p_pin.Nets[i] != nid {
		p_pin.Nets = slices.Insert(p_pin.Nets, i, nid)
		return true
	}

	return false
}

/*
Attach net `nid` to pin `pid`. If either ids are invalid, or the connection
already exists, return false. If a connection is one-sided, repair the
connection and return true.

O(q) for q average connections on a net.
*/
func (nc LCNetController[S, T]) Attach(nid Label, pid Label) bool {

	p_pin, p_net := nc.pinlist.Get(pid), nc.netlist.Get(nid)
	if p_pin == nil || p_net == nil {
		return false
	}

	net_valid := nc.attach_net(nid, pid, p_net)
	pin_valid := nc.attach_pin(nid, pid, p_pin)

	return net_valid || pin_valid
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
			p_source.Pins,
			p_source.Tid,
			p_source.State,
		}
		(*nc.netlist)[nid1] = new_target
	} else {
		//Otherwise, merge the pin arrays, use target for the tid and state
		//Since net pins are sorted, the merge is O(q) with this:
		pins := make([]Label, len(p_target.Pins)+len(p_source.Pins))
		i, j := 0, 0
		for p := range pins {
			//push target
			if i < len(p_target.Pins) && p_target.Pins[i] < p_source.Pins[j] {
				pins[p] = p_target.Pins[i]
				i++
				//remove duplicates
			} else if i < len(p_target.Pins) && j < len(p_source.Pins) && p_target.Pins[i] == p_source.Pins[j] {
				pins[p] = p_target.Pins[i]
				pins = pins[:len(pins)-1]
				i++
				j++
				//push source
			} else if j < len(p_source.Pins) {
				pins[p] = p_source.Pins[j]
				j++
			}
		}

		//set pins, but keep the target's tid and state
		p_target.Pins = pins
	}

	//retarget all pins pointing at source to point at target O(q log q)
	for _, pid := range p_target.Pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		//binary search O(log q)
		i := sort.Search(len(p_pin.Nets), func(i int) bool { return p_pin.Nets[i] >= nid2 })
		if i < 0 || i >= len(p_pin.Nets) {
			continue
		}

		//retarget the pin to point at target
		p_pin.Nets[i] = nid1
	}

	//remove the source *without* removing pins, since they have been retargeted
	nc.netlist.Remove(nid2, LNet[S, T]{})

	return true
}

/*
AddNet a net to the circuit. Copies the input into the circuit, sorting and
validating pins.

If no valid pins are on net.pins, the function fails and returns LABEL_EMPTY

O(q^2) for q average connections
*/
func (nc LCNetController[S, T]) AddNet(net LNet[S, T]) Label {

	//new_net to copy net into. Helps with checking pins
	new_net := LNet[S, T]{
		Pins:  []Label{},
		Tid:   net.Tid,
		State: net.State,
	}
	nid := nc.netlist.Add(new_net, 0)

	//attempt to attach one (1) pin from `net` to `nid`
	//If this fails, don't add the net and return `LABEL_EMPTY`
	valid := -1
	for i, pid := range net.Pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin != nil {
			nc.attach_pin(nid, pid, p_pin)
			(*nc.netlist)[nid].Pins = []Label{pid}
			valid = i
			break
		}
	}

	if valid == -1 {
		var zero LNet[S, T]
		nc.netlist.Remove(nid, zero)
		return LABEL_EMPTY
	}

	//now we know at least one pin was valid, start from the next pin and try to
	//attach the rest. Any valid pins will be sorted by the Attach funciton.
	for _, pid := range net.Pins[valid+1:] {
		nc.Attach(nid, pid)
	}

	return nid
}

/*
Remove a net from the circuit. Disconnects all pins in `net`. If `empty_only`
is true, pins are not removed from nets with a valid `tid` (e.g. gates).

O(nq^2) for n input pins and q average connections
*/
func (nc LCNetController[S, T]) RemoveNet(net LNet[S, T], empty_only bool) bool {

	valid := false
	for _, pid := range net.Pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		nets := make([]Label, len(p_pin.Nets))
		copy(nets, p_pin.Nets)

		for _, nid := range nets {
			if empty_only {
				p_net := nc.netlist.Get(nid)
				if p_net == nil || p_net.Tid != LABEL_EMPTY {
					continue
				}
			}
			valid = nc.Detach(nid, pid) || valid
		}
	}

	return valid
}

func (nc LCNetController[S, T]) ListNets() []Label {
	result := []Label{}
	for nid, net := range *nc.netlist {
		if !net.IsEmpty() {
			result = append(result, Label(nid))
		}
	}

	return result
}

/*
If you thought the previous algorithms were egregious and poorly justified,
get ready! This is a proto-algorithm I came up with for adding/removing wires.

It turns out that splitting nets along a wire requires full knowledge of where
the actual wires are, so I'm going to make an LCWireController that stores the
adjacency list separately. But this was a solid attempt.

Given `net`, "remove" it.

When `net` is removed, for each `pid` in `net`:
  - Remove all `pid2 != pid` in `net` from every `nid` on pid
  - If a pin is fully disconnected, create a new net that belongs to it
  - New nets use the input net's `tid` and `state`.

I call it a "crumble", named after how the process looks when drawn out on paper

Problems with this algorithm:
  - When splitting within a single net, this will create a bunch of loose pins
    instead of "splitting" the pin along a wire.
*/
func (nc LCNetController[S, T]) Crumble(net LNet[S, T]) {
	for _, pid := range net.Pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		for _, pid2 := range net.Pins {
			if pid == pid2 {
				continue
			}

			p_pin2 := nc.pinlist.Get(pid2)
			if p_pin2 == nil {
				continue
			}

			//safe to iterate while detaching since p_pin's network should stay the
			//same
			for _, nid := range p_pin.Nets {
				nc.Detach(nid, pid2)
			}
		}
	}

	for _, pid := range net.Pins {
		p_pin := nc.pinlist.Get(pid)
		if p_pin == nil {
			continue
		}

		if len(p_pin.Nets) == 0 {
			nc.AddNet(LNet[S, T]{
				Pins:  []Label{pid},
				Tid:   net.Tid,
				State: net.State,
			})
		}
	}
}
