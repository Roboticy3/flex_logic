class_name FlexGate extends FlexNet
## A gate with knowledge of its pins. Calculates new values for the pins and
## causes events accordingly.
## Also a likely source of delay, as opposed to pins or wires.

@export var pins:Array[FlexPin]
