package Lcircuit

/*
a pin on a specific component in a circuit
*/
type lpin struct {
	component int64
	pin       int64
}

type lnet[S Lstate] struct {
	pins  []lpin
	state S
}

func (net lnet[S]) IsEmpty() bool {
	return len(net.pins) == 0
}

/*
a component in the circuit has a gate and connections to nets

nets is ordered to match `gtype.pinout`
*/
type lcomponent[S Lstate, T Ltime] struct {
	gtype Lgate[S, T]
	nets  []int64
}

func (comp lcomponent[S, T]) IsEmpty() bool {
	return comp.gtype.name == -1
}

type Lcircuit[S Lstate, T Ltime] struct {
	nets_to_pins  Llabeling[lnet[S]]
	gates_to_nets Llabeling[lcomponent[S, T]]
}

type Lgate_v[S Lstate, T Ltime] struct {
	*Lcircuit[S, T]
}

/*
 Add a gate of type `gtype` to the circuit with no connections, returning the
 gate id

 O(n + m) for n gates and m nets.
*/
func (gview Lgate_v[S, T]) AddGate(gtype Lgate[S, T]) int64 {
	//With no connections, each pin in a gate will induce a separate net on the
	//	circuit.

	// Add the gate component
	gid := gview.gates_to_nets.Add(lcomponent[S, T]{
		gtype: gtype,
		nets:  make([]int64, len(gtype.pinout)),
	}, 0)

	// Sweep through nets to add, binding each one to their respective pin
	nid := int64(0)
	var zero S
	for i := range gtype.pinout {
		nid = gview.nets_to_pins.Add(lnet[S]{
			pins: []lpin{{
				component: gid,
				pin:       int64(i),
			}},
			state: zero,
		}, nid)
		gview.gates_to_nets[gid].nets[i] = nid
	}

	return gid
}

/*
 Remove a gate with id `gid` from the circuit

 O(p) for p pins on the gate.
*/
func (gview Lgate_v[S, T]) RemoveGate(gid int64) bool {
	if gview.gates_to_nets.Get(gid) == nil {
		return false
	}

	for _, nid := range gview.gates_to_nets[gid].nets {
		gview.nets_to_pins.Remove(nid)
	}
	gview.gates_to_nets.Remove(gid)

	return true
}

/*
 Get the type of the gate with id `gid`

 If no gate is found, a gate type passing IsEmpty is returned.
 Breaks if gid is out of range.

 O(1)
*/
func (gview Lgate_v[S, T]) GetGate(gid int64) Lgate[S, T] {
	return gview.gates_to_nets[gid].gtype
}

/*
 List gates currently in the circuit by `gid`

 O(n) for n gates.
*/
func (gview Lgate_v[S, T]) ListGates() []int64 {
	var gids []int64
	for i := 0; i < gview.gates_to_nets.Len(); i++ {
		if !gview.gates_to_nets[i].IsEmpty() {
			gids = append(gids, int64(0))
		}
	}
	return gids
}
