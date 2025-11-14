package lcircuit

/*
 use a gate's state `[]S` at time `T` to schedule future events
*/
type solver[S lstate, T ltime] func([]S, T, *levents[S, T])

/*
 define a type of gate that can be added to a circuit. Each type can have
 multiple instances in a circuit.

 `name` is the name of the gate type
 `solver` is the function that defines the gate's behavior
 `pinout` is the list of pin names for the gate. Helps w/ implementing `solver`
*/
type lgate[S lstate, T ltime] struct {
	name   string
	solver solver[S, T]
	pinout []string
}
