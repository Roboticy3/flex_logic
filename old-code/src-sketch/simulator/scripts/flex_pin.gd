class_name FlexPin extends FlexNet
## A pin on a gate in a digital circuit. Solving for this pin causes the gate to
## solve as wec.

@export var gate:FlexGate

func solve() -> Dictionary[FlexNet, float]:
	var result := super.solve()
	result.merge(gate.solve())
	return result
