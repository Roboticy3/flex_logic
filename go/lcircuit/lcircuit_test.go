package lcircuit

import (
	"testing"
)

func TestAddGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gtypes = testGates

	result := gview.AddGate("AND")

	// Verify that:
	//	1. The returned gate ID is correct
	//	2. A new component is added
	// 	3. An appropriate number of nets was added
	if result != 0 {
		t.Errorf("Expected gate ID to be 0, got %d", result)
	}
	if len(gview.gates_to_nets) != 1 {
		t.Errorf("Expected 1 component in gates_to_nets, got %d", len(gview.gates_to_nets))
	}
	if len(gview.nets_to_pins) != 3 {
		t.Errorf("Expected 3 nets in nets_to_pins, got %d", len(gview.nets_to_pins))
	}
}

func TestAddGateInvalid(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gtypes = testGates

	result := gview.AddGate("OJSDGFDOAFJOEWJF NOT A REAL GATE aoshifdashifuodsoashuifd PLEASE DO NOT NAME GATES LIKE THIS :))))) 🗿")

	if result != -1 {
		t.Errorf("Expected gate ID to be -1 (error) for invalid gate, got %d", result)
	}
	if len(gview.gates_to_nets) != 0 {
		t.Errorf("Expected 0 components in gates_to_nets for invalid gate, got %d", len(gview.gates_to_nets))
	}
	if len(gview.nets_to_pins) != 0 {
		t.Errorf("Expected 0 nets in nets_to_pins for invalid gate, got %d", len(gview.nets_to_pins))
	}
}

func TestAddMultipleGates(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gtypes = testGates

	results := []Label{
		gview.AddGate("AND"),
		gview.AddGate("NOT"),
		gview.AddGate("LATCH"),
	}

	if results[0] != 0 || results[1] != 1 || results[2] != 2 {
		t.Errorf("Expected gate IDs to be 0, 1, 2; got %v", results)
	}
	if len(gview.gates_to_nets) != 3 {
		t.Errorf("Expected 3 slots in gates_to_nets, got %d", len(gview.gates_to_nets))
	}
	if len(gview.nets_to_pins) != 9 { //AND has 3, NOT 2, LATCH 4. No merges = 9 nets
		t.Errorf("Expected 8 slots in nets_to_pins, got %d", len(gview.nets_to_pins))
	}
}

func TestAddRemoveGates(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gtypes = testGates

	gview.AddGate("AND")
	gview.AddGate("NOT")
	gview.AddGate("LATCH")

	gview.RemoveGate(1)

	gates := gview.ListGates()

	if len(gates) != 2 || gates[0] != 0 || gates[1] != 2 {
		t.Errorf("Expected 2 gates with ids 0 (AND) and 2 (LATCH), but found %v", gates)
	}
	if len(gview.gates_to_nets) != 3 {
		t.Errorf("Expected 3 slots in gates_to_nets, got %d", len(gview.gates_to_nets))
	}
	if len(gview.nets_to_pins) != 9 {
		t.Errorf("Expected 9 slots in nets_to_pins, got %d", len(gview.nets_to_pins))
	}
}

func TestFillInRemovedGate(t *testing.T) {
	circuit := &LCircuit[int, int]{}
	gview := LCGateController[int, int]{circuit}
	gview.gtypes = testGates

	gview.AddGate("AND")
	gview.AddGate("NOT")
	gview.AddGate("LATCH")

	gview.RemoveGate(1)

	gview.AddGate("LATCH")

	gates := gview.ListGates()

	if len(gates) != 3 || gates[0] != 0 || gates[1] != 1 || gates[2] != 2 {
		t.Errorf("Expected 3 gates with ids 0 (AND), 1 (second LATCH), and 2 (first LATCH), but found %v", gates)
	}
	if len(gview.gates_to_nets) != 3 {
		t.Errorf("Expected 3 slots in gates_to_nets, got %d", len(gview.gates_to_nets))
	}
	if len(gview.nets_to_pins) != 11 {
		t.Errorf("Expected 11 slots in nets_to_pins, got %d", len(gview.nets_to_pins))
	}
}
