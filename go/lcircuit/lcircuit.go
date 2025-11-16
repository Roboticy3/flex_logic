package lcircuit

import "fmt"

/*
a pin on a specific component in a circuit
*/
type lpin struct {
	component int
	pin       int
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
type lcomponent struct {
	tid  int
	nets []int
}

func (comp lcomponent) IsEmpty() bool {
	return comp.tid == -1
}

type Lcircuit[S Lstate, T Ltime] struct {
	nets_to_pins  Llabeling[lnet[S]]
	gates_to_nets Llabeling[lcomponent]
	gtypes        []Lgate[S, T]
}

type Lgate_v[S Lstate, T Ltime] struct {
	*Lcircuit[S, T]
}

/*
 Add a gate of type `gtype` to the circuit with no connections, returning the
 gate id

 O(n + m) for n gates and m nets.
*/
func (gview Lgate_v[S, T]) AddGate(gname string) int {
	//With no connections, each pin in a gate will induce a separate net on the
	//	circuit.

	fmt.Printf("Adding gate %s\n", gname)
	tid := -1
	for i, gt := range gview.gtypes {
		if gt.name == gname {
			tid = i
			break
		}
	}

	if tid == -1 {
		return -1 // Gate type not found
	}

	// Add the gate component
	gid := gview.gates_to_nets.Add(lcomponent{
		tid:  tid,
		nets: make([]int, len(gview.gtypes[tid].pinout)),
	}, 0)

	// Sweep through nets to add, binding each one to their respective pin
	var zero S
	for i := range gview.gtypes[tid].pinout {
		nid := gview.nets_to_pins.Add(lnet[S]{
			pins: []lpin{{
				component: gid,
				pin:       i,
			}},
			state: zero,
		}, 0)
		fmt.Printf("\tadded pin at nid %d\n", nid)
		gview.gates_to_nets[gid].nets[i] = nid
	}

	return gid
}

/*
 Remove a gate with id `gid` from the circuit. Returns the found gate type if
 successful, or -1

 O(p) for p pins on the gate.
*/
func (gview Lgate_v[S, T]) RemoveGate(gid int) int {

	if gview.gates_to_nets.Get(gid) == nil {
		return -1
	}

	fmt.Printf("Removing gate id %d\n", gid)

	empty_net := lnet[S]{}
	empty_component := lcomponent{-1, []int{}}

	for _, nid := range gview.gates_to_nets[gid].nets {
		fmt.Printf("\tvoiding pin %d\n", nid)
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
func (gview Lgate_v[S, T]) GetGateName(gid int) string {
	return gview.gtypes[gview.gates_to_nets[gid].tid].name
}

/*
 List gates currently in the circuit by `gid`

 O(n) for n gates.
*/
func (gview Lgate_v[S, T]) ListGates() []int {
	var gids []int
	for i := 0; i < gview.gates_to_nets.Len(); i++ {
		if !(gview.gates_to_nets[i].IsEmpty()) {
			gids = append(gids, i)
		}
	}
	return gids
}
