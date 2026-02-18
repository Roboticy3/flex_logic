github.com/Roboticy3/flex_logic_cpp_2

I've been struggling to get audio processing into a good state for months now. Godot is complicated and transparent enough to get lost in incrementalism: doing many small things that look good on their own but don't sum to anything productive. I need to put a more concentrated effort into designing the translation layer between the custom simulation module- which is finished! -and Godot's audio engine.

The current translation layer is technically functional, but very scattered and buggy. Here's what that looks like:

`TapSim`:
 - Implementation of fsim and all the connected components.
`TapLiveSwitch`:
 - Holds a set of pids and a simulator, and has the ability to go "live" iff all the pids are valid.
 - Also handles mutex locking with the simulator to make it editable by scripts without crashing the audio engine.
 - Multiple copies share (!!) a simulator reference. Used by AudioEffects to input and generate audio.
`AudioEffectTapIn`:
 - Use a live switch to pass events into the simulator.
 - Tracks its own clock (!!)
`AudioEffectTapOut`
 - Use a live switch to execute the simulator and read events out as audio
 - Also tracks its own clock (!!)
 - Doesn't process unless an "Ignition" sound source is passed through the bus in-engine.
 - Stores a reference simulator for evaluating the correctness of a circuit.

The messy work of actually configuring these busses to work correctly is all handled through scripting. The biggest issue is the multiple ownerships as a result of different hotfixes accumulating. Technical debt basically.

Here's what I want:

~~`TapSim`~~ -> `TapCircuit`
 - TapSim still owns its simulation functions, but this name is more accurate because the simulation is going to be *driven* from somewhere else. The event propagation methods are more graph-traversal methods. TapSim would also be responsible for implementing things like detecting connectedness in the future.
 - TapCircuit will not own its tick rate.
`TapSimulator`
 - An `AudioStream` with an editor interface of a Dictionary mapping from AudioStreams to pids.
 - Consolidating `TapSimLiveSwitch`, `AudioEffectTapIn`, and `AudioEffectTapOut` will actually make things simpler in this case, since all together this will centralize the circuit execution model. One `mix` method can do what took 3 classes previously.
	 - Because of how audio stream composition and instances work, I would probably need to allocate my own sample buffers for each input to pass into the simulation. Just look out for that.
 - When the audio stream starts playing, invalid pids are ignored instead of pausing the simulation outright.
 - A flag controls whether output is directly audible. This is to get around a limitation of `AudioEffectTapIn` later.
 - The sample-to-sample tick rate of the simulation is now stored here.
 - Also allocates output buffers.
`TapProbe`
 - Reads an output buffer from a `TapSimulator` and plays it from a unique audio stream, which can be positioned in the world.
 - Eventually, I want this to replace raw circuit playback. But the raw can still be positioned. This is a bonus feature.

Proposed implementation steps:
1. Rename TapSim with the existing architecture.
2. Create new test cases.
3. Decide how to approach stream composition with the TapSimulator (is using AudioStreamSynchronized a worthwhile investment?)
4. Implement the input side of TapSimulator and write tests.
	1. Print out the event count on the simulator queue
	2. Sample some specific events.
5. Verify the tests
6. Implement the audible side of TapSimulator and write tests
	1. Event count should stabilize
	2. Add new gate type different from mixer to create a more notable effect
	3. Add mixer variable to slide between summed inputs and simulator output.
	4. Reference simulator should still correctly determine the accuracy of the simulator.
	5. Of course, output should be audible as the sum of outputs.
That would cover all the core features aside from `TapProbe`. The reason I'm not sure how to tackle probing is because, having seen internals of how sound is passed around in Godot, I'm not confident I could do my own stream-to-stream data passage correctly.

Test cases:
1. Custom input streams.
	1. Since I've been looking into custom input streams, I have a sine-wave generator that could be converted into a triangle wave generator easily. Layering these as two inputs would sound good.
	2. For combining them, I suggest a multiplication operator. This would be very sensitive to phase synchronization, and would have a sound unique from mixing. My original tests are with a mixer gate, which is hard to distinguish from other behaviors that emerge from audio programming and just sound like mixing.
