An improved version of [[Naive architectural overview]]. Some things to keep in mind:
1. Types of methods and fields are suggestions. When I say `const T &`, I usually mean "constant (shared) reference". When I say `T[n]`, I mean "n of T", but not necessarily a stack-allocated array (for example, it may be more efficient to store data column-wise). When I say `solver`, that type is not fully resolved. The relations should be cleared up in the UML and the types should be cleared up only when actual code starts getting written down.
2. This architecture is generalized, so the most difficult kind of mistake to spot will be mistakes that make assumptions that are usually valid, but can't be held here. 
3. The architecture is generalized, so optimizations pretty much have to be pushed back to the implementation, which I guess is good?
# Global
## sn_id
A string that is unique under a `sn_id_manager`. StringName is the equivalent primitive type in Godot. In a good implementation, this collapses to a pointer or integer. In the trivial id space case, `sn_id_manager` generates a labeling of the whole numbers ("A", "B", "C"...). This is one step up from enumerators because each key inherently has a (hopefully readable) text description, and some managers can allow for elements to be declared dynamically at runtime using hashing.
1. String to id. If the string is not previously registered, this method will register it.
	1. `const sn_id resolve(string s)`
2. Convert an id back, for display purposes.
	1. `const string &unpack(sn_id id)`

Here are three possible implementations, which could be used as needed:
## sn_idt
`sn_id` with the trivial id layout. Unregistered strings cannot be added, and `resolve` will return an empty id if they are passed. In exchange, all ids have a trivial integer counterpart. 
## sn_idg
`sn_id` with a generalized trivial id layout.
## sn_idh
`sn_id` with a hybrid layout. Most ids will use a trivial id system, but generalized ids can also be added. 

Users of sn_id may need to take a generic implementing sn_id instead, to make the specification dependent on the generics of an interface like other optimization-related choices.
## I/O
Save a load files into a simulator, generating a circuit and routes as needed.
1. Load a simulation from a file path, adding it to the running program. The file may contain circuit data, stream data, and simulation state.
	1. `int load(string path)`.
2. Save a simulation from a file path. If `checkpoint` is true, the simulation's state is encoded along with the active circuit and routes.
	1. `int save(string path, bool checkpoint)`.
# fcircuit module
Each view manages its own association between objects and labels, taking a lot of load off of the circuit and expanding its capability.
## fcgate_c\<S,T>
Add, remove, get and list gates by their label. Added and removed gates induce pins to be added and removed, and mapped to the gates such that the simulator can pass relevant nets to each solver.
1. Add a gate to the circuit under `label`. Can fail if `label` is already in this view. If `label` is not provided, a trivial label is assigned.
	1. `int add_gate(gate<S,T> &g, option<const sn_id> label)`
2. Remove a gate under `label`.
	1. `int remove_gate(const sn_id label)`. 
3. Return a gate under `label` if it exists.
	1. `const option<gate<S,T> &> get_gate(const sn_id label)`. 
4. Returns all labels currently assigned to gates.
	1. `const vector<const sn_id> &list_gates()`. 
## fcpin_c
Add, remove, get and list pins by their label. Each pin belongs to at most one gate. If the circuit is only edited through these controllers, each pin will belong to exactly 1 gate and at most 1 net.
1. Get connections associated with a pin at `label`
	1. `const vector<const net &> &get_connections(const sn_id label)`
2. List all pins.
	1. `const vector<const sn_id> list_pins()`
## fcwire_c
Add, remove, get and list wires by their label and pin endpoints. Wires belonging to the same cluster of pins are consolidated into the same net and the net is assigned an id trivially.
1. Add a wire to the circuit. The full id of a wire is a combination of a label and two pin labels.
	1. `int add_wire(const pair<const sn_id, const sn_id> endpoints, option<const sn_id> label)`, 
2. Remove a wire, either by its label or by a start and end point.
	1. `int remove_wire(const pair<const sn_id, const sn_id> endpoints)`,
	2. `int remove_net(const sn_id label)`. 
3. Get a wire id by its endpoints.
	1. `option<const sn_id> get_wire_id(const pair<const sn_id, const sn_id> endpoints)`
4. Get a wire's endpoints by its id.
	1. `const pair<const sn_id, const sn_id> get_wire_endpoints(const sn_id label)`
5. List wires by endpoints.
	1. `const vector<const pair<const sn_id, const sn_id>> &list_wire_endpoints()`. 
6. List wires by ids
	1. `const vector<const sn_id> &list_wire_ids()`
## fcnet_c
Add, remove, get and list nets by their label and endpoints. A valid circuit can always have its wires displayed, and a net may not correspond to wires directly, so the nets cannot be edited directly.
7. Get a net.
	1. `const const net<S,T> &get_net(sn_id label)`
8. List nets.
	1. `const vector<const sn_id> &list_net()`. 
## fcircuit
Associate gate/net pinout indices to pins. All other connections are secondary to simulation and thus are managed exclusively by their controllers.
1. `map<const sn_id, vector<pair<const sn_id, int>>> pin_connections`
Associate gate/net ids to gate types.
2. `map<const sn_id, const net &> gates`
Associate gates to pins.
3. `map<const sn_id, vector<vector<const sn_id>>> gate_connections`
## net\<S,T>
A net has some default transformation of a state that drives it, typically with very low or zero delay. It has no pinout configuration or name, but does have a solver.
1. `const solver<S,T> solve`
2. The solver takes a list of drivers and returns an `fcevents` queue.
If a `net_manager` exists, it can associate `sn_id` with solvers. `net<S,T>` has an empty id, while `gate<S,T>`, an inheriting interface, must come with a unique, legible id.
## gate\<S,T>
Inherits/composes `net`. Comes with a unique, legible id.
1. `const sn_id name`.
2. `const solver<S,T> solver`. Overrides net's solver.
# fcstream module
Route data structure, the actual Route requirements are satisfied in the fsim module.
## fcstream\<S, T>
An input or output stream that can be routed to an `fcircuit` by an `fsim`. An iterable that can iterate over changes in value or over its value on a certain step.
1. The step size of the stream. This is the step size the last time the stream was written to by a simulation, and is only used for playback.
	1. `T get_step_size()`
2. Loads a stream from a path. Acceptable types: See [[File Types Research]]
	1. `fcstream(string path)`. 
3. Loads a stream from iterable samples and a step size. 
	1. `fcstream(iterable<S> block, T step)`. 
4. Get the next value in the iterable.
	1. `S next()`
## fcbank\<S,T,n>
A collection of streams for the purposes of visualization and serialization. For example: one fcstream may measure one bit of an audio stream, but 16 together are used to save the output to a wav file.
1. The collection of contained streams.
	1. `fcstream<S,T> streams[n]`
2. Save. Determines the saving algorithm from the file extension in `path`. If the extension is not supported, an error is returned.
	1. `sn_id save(string path)`
3. Load. Determines the loading algorithm from the file extension. If the extension is not supported, or the file is corrupt, and error is returned.
	1. `sn_id load(string path)`
4. Error codes for saving and loading.
	1. `sn_id_manager errors`.
## fcevents\<S, T>
Circuit events. A priority queue of `event<S>` ordered by minimum `delay`. Inherits fcstream
1. View the next event
	1. `event<S, T> get_min_after(T time)`
2. Merge an event queue. If `time` is provided, the delay of every event is offset by it, allowing for a time-series append.
	1. `void merge(const fcevents<S, T> &with, T time)`. 
## event\<S, T>
An event on a *pin* in a circuit.
1. `const T time`. 
	1. When produced by a solver, this time is "local" in that it is a delay from when the solver was called. When events are passed into `fcevents::merge`, the current time is applied as an offset, and these times become "global". 
	2. It may be more effective to require solvers to return globally timed events, but then the gate may require more awareness of the circuit than it needs. Technically, this choice can depend on S and T.
2. `const sn_id &drive`
3. `const S value`
# fsim module
Satisfy the Route, Simulator, and Scoring requirements.

The job of `fsim` is to take a circuit and streams and actually run a simulation. In the previous implementation, it's job was to pack the memory layout of the relevant objects closely, so the simulation can retain locality. However, this neglects its most basic purposes of managing the simulation loop, offering debugging interfaces, and writing to outputs.

This is also the most important part to test for efficiency. 

This simulator interface relies on the assumption that fcircuit stores its pin ids internally. This is a fine assumption to make because the circuit will have to at least know its pins' relationships somehow to be a complete graph data structure, so defining the pins and their connections explicitly will always be the simplest way of storing them.
### Setup
1. Set the circuit. If the circuit is the same as the previous circuit, this will force the simulation to acknowledge any changes it might have missed.
	1. `void update_circuit(const fcircuit<S,T> &circuit)`
2. Get the circuit.
	1. `option<const fcircuit<S,T> &> get_circuit()`
3. Add an input stream. If a route is already connected, the route will be driven multiple times on each cycle in the order that the streams were added. If a stream is routed multiple times it will be routed to multiple nets. Each combination of stream and route is unique.
	1. `void add_route_in(const pair<fcstream<S,T> &stream, const pair<const sn_id, const sn_id> route>)`.
4. Remove an input stream. Returns true if the stream was found.
	1. `int remove_route_in(const pair<fcstream<S,T> &stream, const pair<const sn_id, const sn_id> route>)`
5. Get the existing input streams and routes.
	1. `const vector<const pair<fcstream<S,T>, const sn_id>> get_routes_in()`
6. Get the existing routes for a given stream.
	1. `const vector<const sn_id> get_routes_in(const fcstream<S,T> &stream)`
7. Get the existing input streams for a given route.
	1. `const vector<const fcstream<S,T> &> get_streams_in(const pair<const sn_id, const sn_id> route)`
8. Repeat input functions, but for output streams instead.
9. Repeat input functions, but for desired output streams instead.
10. Set the step size for the simulation. All streams are assumed to update at this rate.
	1. `T step_size`
### Simulation
10. Get the current event queue of the simulation.
	1. `const fcevents<S,T> &get_queue()`
	2. All input routes are added to the event queue.
11. Get the current time from start in the simulation. This is equivalent to the current time before the last event propagation, plus the wait time of the event queue.
	1. `const T get_time()`
12. Set the current time of the simulation.
	1. `int scan_to(T time)`
13. Step through one event propagation, keeping the next set of events in memory. Returns false if the step size is reached. If called when no events are in the stack, this is equivalent to `reset`
	1. `int substep()`
	2. When step size is reached, the outputs are written to and the next sample of each of the inputs is added to the event queue with no delay.
	3. Resolve each event to a gate and pin, then update the gate's drivers, and call the gate's solver on the drivers to generate more events.
	4. If the event is part of a route out, append it to the route.
14. Step through one timestep's worth of events. Equivalent to calling substep until it returns false. Returns false if there are no events left in the system, indicating the simulator is done.
	1. `int step()`
15. Get the current score of the simulation. Up to the current `get_time()`, that's the ratio of matching elements between the desired streams and the matching output streams to non-matching elements.
	1. `double score()`
16. Drivers
	1. `map<const sn_id, map<int, S>> drivers`
17. Events. If fcevents implements a `get_min_after(T time)` method instead of pop, this list can be copied selectively to generate traces.
	1. `fcevents<S,T> events`