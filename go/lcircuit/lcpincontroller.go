package lcircuit

import (
	c "flex-logic/collections"
	"flex-logic/ltypes"
	"sort"
)

/*
View and edit a circuit via pins
*/
type LCPinController[S ltypes.LState, T ltypes.LTime] struct {
	*LCircuit[S, T]
}

/*
Get list of nets attached to pin `pid`, or empty list if `pid` is invalid.

Not ideal, since pins can be empty without being invalid.
*/
func (pc LCPinController[S, T]) GetNets(pid c.Label) []c.Label {
	p_pin := pc.pinlist.Get(pid)
	if p_pin == nil {
		return []c.Label{}
	}

	nets := make([]c.Label, len(p_pin.Nets))
	copy(nets, p_pin.Nets)
	return nets
}

/*
Add a pin to the circuit.

If the provided `nid` is `c.LABEL_EMPTY`, the pin is added with no associated
net.

If `nid` is a valid label, but doesn't correspond to a valid net, the function
fails and returns `c.LABEL_EMPTY`

If `nid` corresponds to a valid net, the pin is added to the circuit and
connected to the net.

O(q log q) for q average connections on a net.
*/
func (pc LCPinController[S, T]) AddPin(nid c.Label) c.Label {

	if nid == c.LABEL_EMPTY {
		return pc.pinlist.Add(LPin[S, T]{[]c.Label{}, true}, 0)
	}

	p_net := pc.netlist.Get(nid)
	if p_net == nil {
		return c.LABEL_EMPTY
	}

	result := pc.pinlist.Add(LPin[S, T]{[]c.Label{nid}, true}, 0)
	for _, pid := range p_net.Pins {
		if pid == result {
			return result
		}
	}

	p_net.Pins = append(p_net.Pins, result)
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
func (pc LCPinController[S, T]) RemovePin(pid c.Label) bool {

	p_pin := pc.pinlist.Get(pid)
	if p_pin == nil {
		return false
	}

	//find this pin in the connected nets and remove it
	for _, nid := range p_pin.Nets {
		LCNetController[S, T](pc).Detach(nid, pid)
	}

	var zero LPin[S, T]
	pc.pinlist.Remove(pid, zero)

	return true
}

func (pc LCPinController[S, T]) ListPins() []c.Label {
	result := []c.Label{}
	for pid, pin := range *pc.pinlist {
		if !pin.IsEmpty() {
			result = append(result, c.Label(pid))
		}
	}
	return result
}
