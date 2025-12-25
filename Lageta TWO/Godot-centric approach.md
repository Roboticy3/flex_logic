Making a full circuit simulator outside of godot, then integrating it with godot, failed on 3 points. 
1. Adapting a circuit data structure from simulation-ready structure to editor-ready structure, or making one structure that compromises on either or both forms.
2. Integrating an external system into Godot game engine
3. Implementing the simulator
All of these things proved either unreasonably difficult or impossible for me at my current skill level. I took a break for a few days, but I'm not out of tricks. 

This approach to implementing circuit simulation abuses the audio server built-in with Godot. Specifically, it taps audio samples from audio busses as input to ship them off to processing, then outputs them on audio generators.

Godot busses allow you to tap the samples, which can be read out to a cache. Audio is processed similarly to physics, so the sample count per-step should not change drastically. Then, these caches are read into a mixer. The mixer has cache slots for output. This should be set from some kind of performance settings, with a minimum of 16 to allow for level design and a maximum of 128 for more debug pins. Finally, an AudioStreamGenerator reads the cache into an audio stream. This stream can then be positioned in the world and routed however I, or the player, likes. 

I'm still on break, so I'm going to stop working on this now. This ChatGPT archive has C++ examples that should help me get started when I have time https://chatgpt.com/share/694c36e8-ee30-8003-9561-92db585695fb

I don't like using AI for programming. I wish it was less... bubbly? Like, everything you put in, it's like, "yeah! That's a great idea!" In truth, this is probably going to run into some serious nonsense in practice. Some basic fact-checking reveals that, while a tap is theoretically possible, undocumented doesn't even begin to describe how uncharted this territory is. The base class for it is AudioEffectInstance, not AudioEffect. And a very similar implementation for "Tap" exists called "Capture." Though a custom tap does still happen to be valid because capture is meant for script-side management, I don't think chat knows that.
# Further Details
Different parts of the system and their responsibilities:
 - Circuit builder (GDScript/C++)
	 - Ensures good circuit structure
	 - Relies on game components to build and maintain a circuit in real time.
 - Simulator
	 - Reads the circuit structure on a separate schedule by propagating events through it.
	 - Events are tapped from audio busses, or even other sources. 
	 - Events are outputted through AudioStreamGenerators or other sources. Would be nice to have a color visualization for nets as per the [[Idea for integrating nets]] design.
 - Components/Nets
	 - Process events and pass them to other nets via their pins.

This is nice and all, but it completely defeats the purpose of this approach being *simpler* than the last idea. I wanted the built-in Godot audio components to take care of the component step for me, and they just aren't advanced enough to do that without too much customization. I would also have to deal with C++17 generics again to do this properly. However, I have been working within Godot modules instead of Godot addons. Modules are much more flexible and have access to much more powerful primitives. It may be possible. Just remember, minimum viable product.