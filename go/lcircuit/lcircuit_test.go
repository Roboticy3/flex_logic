package lcircuit

import (
	"testing"
)

func TestAddGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gatetypes = testGates

	result := gview.AddGate("AND")

	// Verify that:
	//	1. The returned gate ID is correct
	//	2. A new component is added
	// 	3. An appropriate number of nets was added
	if result != 0 {
		t.Errorf("Expected gate ID to be 0, got %d", result)
	}
	if len(gview.netlist) != 1 {
		t.Errorf("Expected 1 component in netlist, got %d", len(gview.netlist))
	}
	if len(gview.pinlist) != 3 {
		t.Errorf("Expected 3 nets in pinlist, got %d", len(gview.pinlist))
	}
}

func TestAddGateInvalid(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gatetypes = testGates

	result := gview.AddGate("OJSDGFDOAFJOEWJF NOT A REAL GATE aoshifdashifuodsoashuifd PLEASE DO NOT NAME GATES LIKE THIS :))))) 🗿")

	if result != -1 {
		t.Errorf("Expected gate ID to be -1 (error) for invalid gate, got %d", result)
	}
	if len(gview.netlist) != 0 {
		t.Errorf("Expected 0 components in netlist for invalid gate, got %d", len(gview.netlist))
	}
	if len(gview.pinlist) != 0 {
		t.Errorf("Expected 0 nets in pinlist for invalid gate, got %d", len(gview.pinlist))
	}
}

func TestAddMultipleGates(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gatetypes = testGates

	results := []Label{
		gview.AddGate("AND"),
		gview.AddGate("NOT"),
		gview.AddGate("LATCH"),
	}

	if results[0] != 0 || results[1] != 1 || results[2] != 2 {
		t.Errorf("Expected gate IDs to be 0, 1, 2; got %v", results)
	}
	if len(gview.netlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gview.netlist))
	}
	if len(gview.pinlist) != 9 { //AND has 3, NOT 2, LATCH 4. No merges = 9 nets
		t.Errorf("Expected 8 slots in pinlist, got %d", len(gview.pinlist))
	}
}

func TestAddRemoveGates(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	gc.RemoveGate(1)

	gates := gc.ListGateIds()
	pins := pc.ListPins()

	if len(gates) != 2 || gates[0] != 0 || gates[1] != 2 {
		t.Errorf("Expected 2 valid gates with ids 0 (AND) and 2 (LATCH), but found %v", gates)
	}
	if len(pins) != 7 {
		t.Errorf("Expected 5 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gc.netlist))
	}
	if len(gc.pinlist) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %d", len(gc.pinlist))
	}
}

func TestFillInRemovedGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	gc.RemoveGate(1)

	gc.AddGate("LATCH")

	gates := gc.ListGateIds()
	pins := pc.ListPins()

	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 valid gates with ids 0 (AND), 1 (second LATCH), and 2 (first LATCH), but found %v", gates)
	}
	if len(pins) != 11 {
		t.Errorf("Expected 6 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gc.netlist))
	}
	if len(gc.pinlist) != 11 {
		t.Errorf("Expected 11 slots in pinlist, got %d", len(gc.pinlist))
	}
}

func TestAddPinEmptyNet(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	pc := LCPinController[int, int]{circuit}

	pc.AddPin(LABEL_EMPTY)

	if len(pc.pinlist) != 1 {
		t.Errorf("Expected 1 slot in pinlist, got %d", len(pc.pinlist))
	}
	if len(pc.pinlist[0].nets) != 0 {
		t.Errorf("Expected pin 0 to have 0 nets, got %v", pc.pinlist[0].nets)
	}
}

func TestFillInPin(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	pc := LCPinController[int, int]{circuit}

	pc.AddPin(LABEL_EMPTY)
	middle := pc.AddPin(LABEL_EMPTY)
	pc.AddPin(LABEL_EMPTY)

	pc.RemovePin(middle)

	pc.AddPin(LABEL_EMPTY)
	pins := pc.ListPins()

	if len(pins) != 3 || pins[0] != 0 || pins[1] != 1 || pins[2] != 2 {
		t.Errorf("Expected 2 gates with ids 0 (AND) and 2 (LATCH), but found %v", pins)
	}
	if len(pc.pinlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(pc.netlist))
	}
}

func TestTamperGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	circuit.gatetypes = testGates

	gc := LCGateController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}

	//Add a gate and attempt to tamper with its pinout
	gc.AddGate("AND")
	r1 := pc.AddPin(0)
	r2 := pc.RemovePin(0)

	if r1 == LABEL_EMPTY || !r2 {
		t.Errorf("Expected tampering to succeed. Add: %v, Remove: %v", r1, r2)
	}
	if len(pc.pinlist) != 4 {
		t.Errorf("Expected 4 slots in pinlist, got %v", pc.pinlist)
	}
	if len(gc.netlist) != 1 {
		t.Errorf("Expected 1 slot in netlist, got %v", gc.netlist)
	}
}

func TestMergeFromEmpty(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	result := nc.MergeTwo(0, LABEL_EMPTY)

	gates := gc.ListGateIds()
	pins := pc.ListPins()
	if result {
		t.Errorf("Expected merge to fail, but it succeeded")
	}
	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 valid gates with ids 0 (AND), 1 (NOT), and 2 (LATCH), but found %v", gates)
	}
	if len(pins) != 9 {
		t.Errorf("Expected 5 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.netlist)
	}
	if len(gc.pinlist) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.pinlist)
	}
	if len(gc.netlist[0].pins) != 3 || gc.netlist[0].pins[0] != 0 || gc.netlist[0].pins[1] != 1 || gc.netlist[0].pins[2] != 2 {
		t.Errorf("Expected pins 0, 1, 2 on gate 0, got %v", gc.netlist[0].pins)
	}
}

func TestMergeIntoEmpty(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	result := nc.MergeTwo(LABEL_EMPTY, 0)

	gates := gc.ListGateIds()
	pins := pc.ListPins()
	if result {
		t.Errorf("Expected merge to fail, but it succeeded")
	}
	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 valid gates with ids 0 (AND), 1 (NOT), and 2 (LATCH), but found %v", gates)
	}
	if len(pins) != 9 {
		t.Errorf("Expected 5 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.netlist)
	}
	if len(gc.pinlist) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.pinlist)
	}
	if len(gc.netlist[0].pins) != 3 || gc.netlist[0].pins[0] != 0 || gc.netlist[0].pins[1] != 1 || gc.netlist[0].pins[2] != 2 {
		t.Errorf("Expected pins 0, 1, 2 on gate 0, got %v", gc.netlist[0].pins)
	}
}

func TestMergeTwo(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	nc := LCNetController[int, int]{circuit}
	result := nc.MergeTwo(0, 2)

	gates := gc.ListGateIds()
	pins := pc.ListPins()
	if !result {
		t.Errorf("Expected merge to succeed, but it failed")
	}
	if len(gates) != 2 || gates[0] != 0 || gates[1] != 1 {
		t.Errorf("Expected 2 valid gates with ids 0 (AND) and 1 (NOT), but found %v", gates)
	}
	if len(pins) != 9 {
		t.Errorf("Expected 9 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.netlist)
	}
	if len(gc.pinlist) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.pinlist)
	}
}

func TestMergeMany(t *testing.T) {
	//Test merging many nets, invalid and valid, into one
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	pc := LCPinController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")
	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")
	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")
	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	gc.RemoveGate(6)
	gc.RemoveGate(10)

	gates := gc.ListGateIds()
	pins := pc.ListPins()

	if len(gates) != 10 {
		t.Errorf("Expected 10 valid gates, but found %v", gates)
	}
	if len(pins) != 31 {
		t.Errorf("Expected 31 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 12 {
		t.Errorf("Expected 12 slots in netlist, got %v", gc.netlist)
	}
	if len(gc.pinlist) != 36 {
		t.Errorf("Expected 36 slots in pinlist, got %v", gc.pinlist)
	}

	/*
		Test:
			-1: invalid id at the beginning of the list
			6, 10: ids of removed gates
			others: different gate types
			not included: leaving certain gates alone

		Desired effect:
			0, 4, 5, 7, 8, 11 merged into 0. 10 gates becomes 5. Ids:
				0, 1, 2, 3, 9
			36 pins added, 5 removed. Pins are sustained through merge. Total 31 valid
			12 nets added.
	*/
	nc.Merge([]Label{-1, 0, 4, 5, 6, 7, 8, 10, 11})

	gates = gc.ListGateIds()
	pins = pc.ListPins()
	if len(gates) != 5 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 || gates[3] != 3 || gates[4] != 9 {
		t.Errorf("Expected 5 valid gates with ids 0 (AND), 1 (NOT), 2 (LATCH), 3 (AND), and 9 (AND), but found %v", gates)
	}
	if len(pins) != 31 {
		t.Errorf("Expected 31 valid pins, got %v", pins)
	}
	if len(gc.netlist) != 12 {
		t.Errorf("Expected 12 slots in netlist, got %v", gc.netlist)
	}
	if len(gc.pinlist) != 36 {
		t.Errorf("Expected 36 slots in pinlist, got %v", gc.pinlist)
	}
}

func TestAddNet(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := LNet[int, int]{
		[]Label{8, 4, 2, 0, -1, 345, 4}, -1, -1,
	}
	nid := nc.AddNet(new_net)
	p_net := nc.netlist.Get(nid)

	nets := nc.ListNets()

	if p_net == nil {
		t.Errorf("Must be able to retrieve net %d to continue with test. valid nets: %v", nid, nets)
		return
	}

	/*
		Adding a net should sort the input pins
		Invalid pins should be skipped

		Expected pinout of new net: 0, 2, 4, 8
	*/
	if len(nets) != 4 || nets[0] != 0 || nets[1] != 1 || nets[2] != 2 || nets[3] != 3 {
		t.Errorf("Expected 4 valid nets with contiguous ids, got %v", nets)
	}

	pins := p_net.pins
	if len(pins) != 4 || pins[0] != 0 || pins[1] != 2 || pins[2] != 4 || pins[3] != 8 {
		t.Errorf("Expected net %d to have pins 0, 2, 4, 9, found %v", nid, pins)
	}

}

func TestFailAddNet(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := LNet[int, int]{
		[]Label{-1, -3244, 45, 9}, -1, -1,
	}
	nid := nc.AddNet(new_net)
	p_net := nc.netlist.Get(nid)

	/*
		Adding a net fails when no valid pins are on the net. Since empty nets are
		invalid states, the net should be removed cleanly.
	*/

	if p_net != nil {
		t.Errorf("Expected to fail adding net, found net %d with pins %v", nid, p_net.pins)
		return
	}
}

func TestRemoveNetClean(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := LNet[int, int]{
		[]Label{8, 4, 2, 0, -1, 345, 4}, -1, -1,
	}
	nid := nc.AddNet(new_net)

	/*
		Removing all pins on a net should remove it from the circuit
	*/
	result := nc.RemoveNet(new_net, false)
	nets := nc.ListNets()

	if !result {
		t.Errorf("Expected to remove net %d, but it failed", nid)
	}
	if len(nets) != 3 {
		t.Errorf("Expected 3 valid nets, got %v", nets)
	}
}

func TestRemoveNetDirty(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := LNet[int, int]{
		[]Label{8, 4, 2, 0, -1, 345, 4}, -1, -1,
	}
	nid := nc.AddNet(new_net)

	/*
		Try removing a set of indices not corresponding directly to an existing net.
		This test case made more sense when the remove behavior was more complicated
		, but it still makes me comfortable to have it around.
	*/
	dirty_remove := LNet[int, int]{
		[]Label{0, 1, 2, -1, 9}, -1, -1,
	}
	result := nc.RemoveNet(dirty_remove, false)
	nets := nc.ListNets()

	if !result {
		t.Errorf("Expected to remove net %d, but it failed", nid)
	}
	if len(nets) != 3 {
		t.Errorf("Expected 3 valid nets, got %v", nets)
	}
}

func TestRemoveNetInvalid(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gc := LCGateController[int, int]{circuit}
	nc := LCNetController[int, int]{circuit}
	gc.gatetypes = testGates

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := LNet[int, int]{
		[]Label{8, 4, 2, 0, -1, 345, 4}, -1, -1,
	}
	nc.AddNet(new_net)

	/*
		Try removing only invalid indices. Netlist should be unaffected and result
		should be false
	*/
	invalid_remove := LNet[int, int]{
		[]Label{-1, 9, 2359, -3124}, -1, -1,
	}
	result := nc.RemoveNet(invalid_remove, false)
	nets := nc.ListNets()

	if result {
		t.Errorf("Expected to fail removing %v, instead succeeded.", invalid_remove.pins)
	}
	if len(nets) != 4 {
		t.Errorf("Expected 4 valid nets, got %v", nets)
	}
}
