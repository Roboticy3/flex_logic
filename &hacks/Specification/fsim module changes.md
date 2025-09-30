Much like [[fcircuit module changes]], fsim could benefit from splitting. Adding a controller to access the state inherently solves the problem of implementing a trace. 

The fcircuit changes also allow for more specification of how fsim will pass input to gates.
## event\<S, T>
One event on a *pin*
1. `const T time`
2. `const sn_id &drive`
3. `const S value`
# fsim\<S,T>
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
## fstrace
A complete or partial trace of a simulator. Whenever an event is processed, it is added to the trace. Traces are a type of `fcstream` and can be added to the simulator with `add_route_out`.
1. Events.
	1. `fcevents<S,T> events`