So, how would a simulator like this work? How would it be implemented?

My implementation will be inspired by the IEEE Verilog standard, and seeks to get a balance between realism and usability. 

# Data Structures
## Net
A Net is a Node connected to a list of other Nodes that can be initialized to a set of IEEE Verilog net states. It also references a Sim which will handle it's actual behavior.

As the Sim handles its actual state, the Net's state can change to reflect it. This resolves the issue of the Sim possibly running at a much higher frequency than the refresh rate of the game. The Net state should reflect the "visual" state, and allow for a deferred setting of the value at runtime.

I would really like for the state to be a template argument, since then the interactions between the states specified in the IEEE could be implemented separately from the base class. But Godot C++ and gdscript have terrible support for this, it may be worth looking into other language plugins like Rust or Go.
## Sim
A simulator is a Node that collects a list of Nets and compacts their connections into a simple index map.

Each Net corresponds to a "resolver" method by storing a resolver index. The Sim then stores a list of resolvers. Not sure if I want to do this as classes or methods. A resolver takes an event queue and a pointer to a set of values corresponding to the Sim's managed state of a Net. To resolve a Net, the Sim can call its associated resolver function with an event queue and the Net's managed state as input. The Sim will then provide all the "built-in" simulators provided by the game. The default sim is a Wire. When a wire connects to another, it "drives" it, in a behavior that is well-defined by the IEEE. 
 > It may be worth looking into using a single method with a switch. This method would not only be slightly faster, but also easier to read, in exchange for less flexibility. 

The Sim holds the functions to actually effect the simulator state. From simplest to most complex:
1. `step_once_from`
	1. Create an event queue from a starting Net, 
	2. then push all it's connections to the event queue. 
	3. Then resolve those events. Then stop.
2. `step_path`
	1. Deduce the shortest path between any two Nets.
	2. Call `step_once_from` on each Net in the path.
3. `step`
	1. Resolve all events in a stateful event queue belonging to the Sim, this pushes more events
	2. Repeat until no events remain.

To use the Sim as a full simulator, it should expose a function to manually push events. 
1. `update_single`
	1. Change the state of a gate or an input gate.
	2. (Maybe automatically) push the event of the input gate onto the simulator.
	3. Call `resolve`
2. `update_block`
	1. Send a matrix of inputs to the Sim
	2. For each item in the block, call `update_single`.
The Sim needs a notion of input and output nets for `update_block` to work automatically. A good way to do it would be to give Nets a mapping between "groups" and "id_priority", so the Sim can organize them into groups with contiguous id spaces, which can then be invoked by group from `update_block`. But.. That would also be really hard. Update Block might have to wait. Think about simpler puzzles you could demo without it.
### Delay
The event queue will be in FILO priority by default, but events can have a lower priority to give them a delay. The step and update functions should take a `delta` time describing the max time they have to clear events, so that any delayed events stay on the event queue until the next cac.
## Gate
A Gate is a special kind of Net that creates a set of sub-nets from which it draws its inputs. This linkage allows it to know what order its inputs are stored in, which is important for corresponding with resolvers.
## All together
To build a simulation, devise a method by which gates and wires can be constructed, and added in real time to the Sim. Next, connect a block of inputs and expected outputs to run `update_block` on demand, progressing by a set timestep with high enough frequency to produce audio (may not be possible in scripting). 

Gates and Wires should properly indicate their state.
# GDScript implementation
First, to get a general idea of how the simulation should work, I should make a gdscript implementation. gdscript is slower, but easier to debug. In my experience, it's very hard to get breakpoints set up with GDExtension C++, something I have continually struggled with, and they don't seem to be difficulties that would go away by just changing my plugin language, since they have to do with the complexity of how Godot initializes debugging processes. If I find a solution to that, the scripting implementation may still be relevant, since it's easier in scripting to experiment with different export properties. 

Currently, I have an unfinished C++ implementation. It may be possible to replace it in-place(!!) with a GDScript prototype, as long as I don't open the editor until all the type names are resolved.

Faster experimenting would also help for simplifying certain parts of the simulation that I'm not clear on yet. For example, the state management of Nets, and the way they're properties proxy the Sim, is still pretty unclear. The most similar data structure already existing in godot is [PhysicsServer3D](https://docs.godotengine.org/en/stable/classes/class_physicsserver3d.html) or PhysicsServer2D, which works in a similar way.
# Compiled implementation
This is where the Sim architecture should really shine. All the gates states and relationships should get packed into contiguous blocks of memory, allowing for fast lookups and resolutions. Ideally, the event queue, state array, and connections array are all cached in part at any given time, and the state is a template type.