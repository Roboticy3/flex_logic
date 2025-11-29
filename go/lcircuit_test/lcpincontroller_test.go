package lcircuit_test

import (
	"flex-logic/lcircuit"
	"testing"
)

func TestAddPinEmptyNet(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}

	pid := pc.AddPin(lcircuit.LABEL_EMPTY)

	if len(pc.GetPinlist()) != 1 {
		t.Errorf("Expected 1 slot in pinlist, got %d", len(pc.GetPinlist()))
	}
	if len(pc.GetNets(pid)) != 0 {
		t.Errorf("Expected pin 0 to have 0 nets, got %v", pc.GetNets(pid))
	}
}

func TestFillInPin(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}

	pc.AddPin(lcircuit.LABEL_EMPTY)
	middle := pc.AddPin(lcircuit.LABEL_EMPTY)
	pc.AddPin(lcircuit.LABEL_EMPTY)

	pc.RemovePin(middle)

	pc.AddPin(lcircuit.LABEL_EMPTY)
	pins := pc.ListPins()

	if len(pins) != 3 || pins[0] != 0 || pins[1] != 1 || pins[2] != 2 {
		t.Errorf("Expected 2 gates with ids 0 (AND) and 2 (LATCH), but found %v", pins)
	}
	if len(pc.GetPinlist()) != 3 {
		t.Errorf("Expected 3 slots in netlist, got %d", len(pc.GetNetlist()))
	}
}

func TestTamperGate(t *testing.T) {
	circuit := lcircuit.CreateCircuit[int, int]()
	circuit.SetGateTypes(testGates)

	gc := lcircuit.LCGateController[int, int]{LCircuit: circuit}
	pc := lcircuit.LCPinController[int, int]{LCircuit: circuit}

	//Add a gate and attempt to tamper with its pinout
	gc.AddGate("AND")
	r1 := pc.AddPin(0)
	r2 := pc.RemovePin(0)

	if r1 == lcircuit.LABEL_EMPTY || !r2 {
		t.Errorf("Expected tampering to succeed. Add: %v, Remove: %v", r1, r2)
	}
	if len(pc.GetPinlist()) != 4 {
		t.Errorf("Expected 4 slots in pinlist, got %v", pc.GetPinlist())
	}
	if len(gc.GetNetlist()) != 1 {
		t.Errorf("Expected 1 slot in netlist, got %v", gc.GetNetlist())
	}
}
