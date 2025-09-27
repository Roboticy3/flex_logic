# Simulator
Two constructs provide different "views" of the simulator: the stream and the circuit.

Streams drive the inputs of a circuit, observe its outputs during simulation, and can be observed as a waveform. They have the real life equivalent of testing a circuit board with a volt meter. 

Circuits allow a user to view the structure of the simulator. They have the real life equivalent of observing a layout of chips on a breadboard and moving around the chips and wires.

The simulation itself is like the actual electricity in a circuit. For example, in a circuit following the verilog standard will carry values 0,1,X,Z, or U in each net. 0 and 1 be observed as high and low in a square wave, while illegal values get measured as different flavors of noise.

Note that the clock does not correspond to real time. In fact, it will likely be slower than real time in most cases. The clock describes how much time is passing *inside* the simulation, not how the circuit should synchronize to real life, though it can be used for that purpose.

Importantly, no part of these requirements limits the actual data type passed through the nets. For example, this simulator is a sufficient interface to implement a quantum circuit simulation, as long as the clock can be adjusted to a reasonable step and the waveform interpreted correctly. Each net could carry a 2d complex vector that is transformed by linear operators, and nets could be limited to one input and output due to the no-cloning theorem. 

The most interesting net type for the purposes of this project is a vector type. This is natively supported by VHDL and corresponds nicely to bitwise audio waveform. For example, an 8-bit adder circuit would have the exact effect of combining two 8-bit audio streams, with a carry bit indicating if the audio is peaking.
## Routing/Streams
1. A stream can be routed to another stream directly and "run" on a clock, copying one stream to the other at a set rate.
2. A stream can be routed to a port on a circuit and run on a clock, copying the stream to the port as an input.
3. A port can be routed to a stream in the same way to produce an output.
4. Streams can be saved and loaded. A loaded stream will have its clock speed encoded.
5. As a stream is playing, its value can be observed as a waveform.
6. As a stream is playing into a circuit, the circuit will be simulated.
7. If a circuit simulated from a stream is observed from a stream, that output can also be observed as a waveform
8. In deterministic simulations, a set of input streams and "desired" output streams defines a set of circuits that can take the input streams to some ports and produce the same outputs.
9. In deterministic simulations, all combinations of circuits, inputs, and desired outputs have a proportional "score" out of 1, depending on how well the circuit matches the outputs.
## Circuit
1. A circuit can be created, saved, and loaded from a file
2. An open circuit can have gates added and removed. Each gate can be inspected to describe its port layout and simulation behavior
	1. Question: should gate types be encoded into a circuit, or hard coded? Hard-coding would be good for prototyping until I have more information. 
3. An open circuit can have a connection made between any two ports at a time.
4. Ports can be configured to load into a certain state, which may effect the score of a circuit.
5. Circuit files can be paired with saved inputs and expected outputs into a Level with a score requirement. The circuit does not have a complete score, but rather is a template with input and output gates to be completed in a circuit editor.
## Simulation
1. A circuit can be simulated in timesteps. 
2. If a circuit is simulated from a stream, the timestep will be the stream's clock.
3. Multiple wires on the same port are consolidated into a net. 
4. One net can correspond to one value.
5. Setting the value of a port "drives" the net.
6. Driving a net causes an event. Events propagate through gates and nets with a delay depending on the net/gate type, adding events to the corresponding ports. 
7. Gates typically only drive their output ports, while nets drive all connected ports except for the driving port. 
8. Events are processed in order of delay up to and including the timestep. Higher delays are reserved until the next simulation step.
9. Simulation can be paused at any step, and the state of the circuit retained and inspected.
10. Simulation can be paused at any substep of event propagation. Unprocessed and processed events can be differentiated.
11. Individual steps and substeps can be stepped.
12. Both driving (input) and observing (output) streams can be added and removed during simulation.
