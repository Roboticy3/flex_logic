package lcircuit_test

import (
	c "flex-logic/collections"
	"flex-logic/lcircuit"
	"testing"
)

func TestMergeFromEmpty(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	result := nc.MergeTwo(0, c.LABEL_EMPTY)

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
	if len(gc.GetNetlist()) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.GetNetlist())
	}
	if len(gc.GetPinlist()) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.GetPinlist())
	}
	if len(nc.GetPins(0)) != 3 || nc.GetPins(0)[0] != 0 || nc.GetPins(0)[1] != 1 || nc.GetPins(0)[2] != 2 {
		t.Errorf("Expected pins 0, 1, 2 on gate 0, got %v", nc.GetPins(0))
	}
}

func TestMergeIntoEmpty(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	result := nc.MergeTwo(c.LABEL_EMPTY, 0)

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
	if len(gc.GetNetlist()) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.GetNetlist())
	}
	if len(gc.GetPinlist()) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.GetPinlist())
	}
	if len(nc.GetPins(0)) != 3 || nc.GetPins(0)[0] != 0 || nc.GetPins(0)[1] != 1 || nc.GetPins(0)[2] != 2 {
		t.Errorf("Expected pins 0, 1, 2 on gate 0, got %v", nc.GetPins(0))
	}
}

func TestMergeTwo(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
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
	if len(gc.GetNetlist()) != 3 {
		t.Errorf("Expected 2 slots in netlist, got %v", gc.GetNetlist())
	}
	if len(gc.GetPinlist()) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %v", gc.GetPinlist())
	}
}

func TestMergeMany(t *testing.T) {
	//Test merging many nets, invalid and valid, into one
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

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
	if len(gc.GetNetlist()) != 12 {
		t.Errorf("Expected 12 slots in netlist, got %v", gc.GetNetlist())
	}
	if len(gc.GetPinlist()) != 36 {
		t.Errorf("Expected 36 slots in pinlist, got %v", gc.GetPinlist())
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
	nc.Merge([]c.Label{-1, 0, 4, 5, 6, 7, 8, 10, 11})

	gates = gc.ListGateIds()
	pins = pc.ListPins()
	if len(gates) != 5 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 || gates[3] != 3 || gates[4] != 9 {
		t.Errorf("Expected 5 valid gates with ids 0 (AND), 1 (NOT), 2 (LATCH), 3 (AND), and 9 (AND), but found %v", gates)
	}
	if len(pins) != 31 {
		t.Errorf("Expected 31 valid pins, got %v", pins)
	}
	if len(gc.GetNetlist()) != 12 {
		t.Errorf("Expected 12 slots in netlist, got %v", gc.GetNetlist())
	}
	if len(gc.GetPinlist()) != 36 {
		t.Errorf("Expected 36 slots in pinlist, got %v", gc.GetPinlist())
	}
}

func TestAddNet(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := lcircuit.LNet{
		Pins: []c.Label{8, 4, 2, 0, -1, 345, 4},
		Tid:  -1,
	}
	nid := nc.AddNet(new_net)
	pins := nc.GetPins(nid)
	nets := nc.ListNets()

	if nid == c.LABEL_EMPTY {
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
	if len(pins) != 4 || pins[0] != 0 || pins[1] != 2 || pins[2] != 4 || pins[3] != 8 {
		t.Errorf("Expected net %d to have pins 0, 2, 4, 9, found %v", nid, pins)
	}

}

func TestFailAddNet(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := lcircuit.LNet{
		Pins: []c.Label{8, 4, 2, 0, -1, 345, 4},
		Tid:  -1,
	}
	nid := nc.AddNet(new_net)
	pins := nc.GetPins(nid)

	/*
		Adding a net fails when no valid pins are on the net. Since empty nets are
		invalid states, the net should be removed cleanly.
	*/

	if nid == c.LABEL_EMPTY {
		t.Errorf("Expected to fail adding net, found net %d with pins %v", nid, pins)
		return
	}
}

func TestRemoveNetClean(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := lcircuit.LNet{
		Pins: []c.Label{8, 4, 2, 0, -1, 345, 4},
		Tid:  -1,
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
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := lcircuit.LNet{
		Pins: []c.Label{8, 4, 2, 0, -1, 345, 4},
		Tid:  -1,
	}
	nid := nc.AddNet(new_net)

	/*
		Try removing a set of indices not corresponding directly to an existing net.
		This test case made more sense when the remove behavior was more complicated
		, but it still makes me comfortable to have it around.
	*/
	dirty_remove := lcircuit.LNet{
		Pins: []c.Label{0, 1, 2, -1, 9},
		Tid:  -1,
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
	circuit := lcircuit.CreateCircuit[int, int]()
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("LATCH")

	new_net := lcircuit.LNet{
		Pins: []c.Label{8, 4, 2, 0, -1, 345, 4},
		Tid:  -1,
	}
	nc.AddNet(new_net)

	/*
		Try removing only invalid indices. Netlist should be unaffected and result
		should be false
	*/
	invalid_remove := lcircuit.LNet{
		Pins: []c.Label{-1, -324, 456, 9},
		Tid:  -1,
	}
	result := nc.RemoveNet(invalid_remove, false)
	nets := nc.ListNets()

	if result {
		t.Errorf("Expected to fail removing %v, instead succeeded.", invalid_remove.Pins)
	}
	if len(nets) != 4 {
		t.Errorf("Expected 4 valid nets, got %v", nets)
	}
}

func TestComplexNetwork(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	/*
		When writing TestCrumble, I ran into some cases when AddNet doesn't work
		as expected. Let's see what's wrong.
	*/

	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)

	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{3, 4, 7, 8, 9},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{9, 10, 11, 12},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{2, 5, 6, 9, 10},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{0, 1, 2, 3, 4},
	})

	nets := nc.ListNets()
	if len(nets) != 4 {
		t.Errorf("Expected 5 valid nets, got %v", nets)
	}

	connections := [][]c.Label{
		nc.GetPins(0),
		nc.GetPins(1),
		nc.GetPins(2),
		nc.GetPins(3),
	}

	if len(connections[0]) != 5 || len(connections[1]) != 4 || len(connections[2]) != 5 || len(connections[3]) != 5 {
		t.Errorf("Expected 5, 4, 5 pins, got %v", connections)
	}
}

func TestCrumble(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	nc := lcircuit.LCNetController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	/*
		We need 13 pins for this.
	*/
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)
	pc.AddPin(c.LABEL_EMPTY)

	/*
		the example I drew on paper has the following four nets to start.
		  0: 3, 4, 7, 8, 9
			1: 9, 10, 11, 12
			2: 2, 5, 6, 9, 10
			3: 0, 1, 2, 3, 4

		I then crumble the following net:
			5, 6, 8, 9, 10

		The expected output is then:
			0: 3, 4, 7, 8
			1: 9, 11, 12
			2: 2, 5
			3: 0, 1, 2, 3, 4 (unaffected)
	*/

	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{3, 4, 7, 8, 9},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{9, 10, 11, 12},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{2, 5, 6, 9, 10},
	})
	nc.AddNet(lcircuit.LNet{
		Pins: []c.Label{0, 1, 2, 3, 4},
	})

	nc.Crumble(lcircuit.LNet{
		Pins: []c.Label{5, 6, 8, 9, 10},
	})

	nets := nc.ListNets()
	pins := pc.ListPins()
	if len(nets) != 4 {
		t.Errorf("Expected 4 valid nets, got %v", nets)
	}
	if len(pins) != 13 {
		t.Errorf("Expected 13 valid pins, got %v", pins)
	}

	connections := [][]c.Label{
		nc.GetPins(0),
		nc.GetPins(1),
		nc.GetPins(2),
		nc.GetPins(3),
	}

	if len(connections[0]) != 4 || len(connections[1]) != 3 || len(connections[2]) != 2 || len(connections[3]) != 5 {
		t.Errorf("Expected 4, 3, 2, 5 pins, got %v", connections)
	}

}
