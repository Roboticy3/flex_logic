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
	gview := LCGateController[int, int]{circuit}
	gview.gatetypes = testGates

	gview.AddGate("AND")
	gview.AddGate("NOT")
	gview.AddGate("LATCH")

	gview.RemoveGate(1)

	gates := gview.ListGateIds()

	if len(gates) != 2 || gates[0] != 0 || gates[1] != 2 {
		t.Errorf("Expected 2 gates with ids 0 (AND) and 2 (LATCH), but found %v", gates)
	}
	if len(gview.netlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gview.netlist))
	}
	if len(gview.pinlist) != 9 {
		t.Errorf("Expected 9 slots in pinlist, got %d", len(gview.pinlist))
	}
}

func TestFillInRemovedGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gatetypes = testGates

	gview.AddGate("AND")
	gview.AddGate("NOT")
	gview.AddGate("LATCH")

	gview.RemoveGate(1)

	gview.AddGate("LATCH")

	gates := gview.ListGateIds()

	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 gates with ids 0 (AND), 1 (second LATCH), and 2 (first LATCH), but found %v", gates)
	}
	if len(gview.netlist) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(gview.netlist))
	}
	if len(gview.pinlist) != 11 {
		t.Errorf("Expected 11 slots in pinlist, got %d", len(gview.pinlist))
	}
}
