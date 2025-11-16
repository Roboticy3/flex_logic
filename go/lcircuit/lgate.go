package lcircuit

/*
 use a gate's state `[]S` at time `T` to schedule future events
*/
type Solver[S Lstate, T Ltime] func([]S, T, *levents[S, T])

/*
 define a type of gate that can be added to a circuit. Each type can have
 multiple instances in a circuit.

 `name` is the name of the gate type
 `Solver` is the function that defines the gate's behavior
 `pinout` is the list of pin names for the gate. Helps w/ implementing `Solver`
*/
type Lgate[S Lstate, T Ltime] struct {
	name   string
	Solver Solver[S, T]
	pinout []string
}
