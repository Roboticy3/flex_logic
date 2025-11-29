package lcircuit_test

import (
	"flex-logic/lcircuit"
	"testing"
)

func TestAddGate(t *testing.T) {
	circuit := &lcircuit.LCircuit[int, int]{}
	gview := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	result := gview.AddGate("AND")

	// Verify that:
	//	1. The returned gate ID is correct
	//	2. A new component is added
	// 	3. An appropriate number of nets was added
	if result != 0 {
		t.Errorf("Expected gate ID to be 0, got %d", result)
	}
	if len(gview.GetNetlist()) != 1 {
		t.Errorf("Expected 1 component in netlist, got %d", len(gview.GetNetlist()))
	}
	if len(gview.GetPinlist()) != 3 {
		t.Errorf("Expected 3 nets in pinlist, got %d", len(gview.GetPinlist()))
	}
}

func TestAddGateInvalid(t *testing.T) {
	circuit := &lcircuit.LCircuit[int, int]{}
	gview := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	result := gview.AddGate("OJSDGFDOAFJOEWJF NOT A REAlcircuit.L GATE aoshifdashifuodsoashuifd Plcircuit.LEASE DO NOT NAME GATES lcircuit.LIKE THIS :))))) 🗿")

	if result != -1 {
		t.Errorf("Expected gate ID to be -1 (error) for invalid gate, got %d", result)
	}
	if len(gview.GetNetlist()) != 0 {
		t.Errorf("Expected 0 components in netlist for invalid gate, got %d", len(gview.GetNetlist()))
	}
	if len(gview.GetPinlist()) != 0 {
		t.Errorf("Expected 0 nets in pinlist for invalid gate, got %d", len(gview.GetPinlist()))
	}
}

func TestAddMultipleGates(t *testing.T) {
	circuit := &lcircuit.LCircuit[int, int]{}
	gview := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	results := []lcircuit.Label{
		gview.AddGate("AND"),
		gview.AddGate("NOT"),
		gview.AddGate("lcircuit.LATCH"),
	}

	if results[0] != 0 || results[1] != 1 || results[2] != 2 {
		t.Errorf("Expected gate IDs to be 0, 1, 2; got %v", results)
	}
	if len(gview.GetNetlist()) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gview.GetNetlist()))
	}
	if len(gview.GetPinlist()) != 9 { //AND has 3, NOT 2, lcircuit.LATCH 4. No merges = 9 nets
		t.Errorf("Expected 8 slots in pinlist, got %d", len(gview.GetPinlist()))
	}
}

func TestAddRemoveGates(t *testing.T) {
	circuit := &lcircuit.LCircuit[int, int]{}
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("lcircuit.LATCH")

	gc.RemoveGate(1)

	gates := gc.ListGateIds()
	pins := pc.ListPins()

	if len(gates) != 2 || gates[0] != 0 || gates[1] != 2 {
		t.Errorf("Expected 2 valid gates with ids 0 (AND) and 2 (lcircuit.LATCH), but found %v", gates)
	}
	if len(pins) != 7 {
		t.Errorf("Expected 5 valid pins, got %v", pins)
	}
	if len(gc.GetNetlist()) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gc.GetNetlist()))
	}
	if len(gc.GetPinlist()) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %d", len(gc.GetPinlist()))
	}
}

func TestFillInRemovedGate(t *testing.T) {
	circuit := &lcircuit.LCircuit[int, int]{}
	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}
	circuit.SetGateTypes(testGates)

	gc.AddGate("AND")
	gc.AddGate("NOT")
	gc.AddGate("lcircuit.LATCH")

	gc.RemoveGate(1)

	gc.AddGate("lcircuit.LATCH")

	gates := gc.ListGateIds()
	pins := pc.ListPins()

	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 valid gates with ids 0 (AND), 1 (second lcircuit.LATCH), and 2 (first lcircuit.LATCH), but found %v", gates)
	}
	if len(pins) != 11 {
		t.Errorf("Expected 6 valid pins, got %v", pins)
	}
	if len(gc.GetNetlist()) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gc.GetNetlist()))
	}
	if len(gc.GetPinlist()) != 11 {
		t.Errorf("Expected 11 slots in pinlist, got %d", len(gc.GetPinlist()))
	}
}
