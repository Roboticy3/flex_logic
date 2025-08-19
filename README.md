# Godot 4 Digital Logic simulator
Goal: Make a digital logic simulator for my game. It should be fun and intuitive to play with first and foremost.

Features:
1. Cycle-based simulation of digital logic circuits. More performant than events for gdscript.
2. Integer-based logic to allow for seamless integration of buses. 
3. Live editing of circuits along with some primitives for the user interface that I can build off of.

Cool features:
1. Piping through audio streams. Godot theoretically has support for this out of the box. The integer arithmetic will help for this, but not for the reason I thought. I thought I would be able to supply integers to the audio server bitwise, but Godot uses an abstraction called [AudioFrame](https://github.com/godotengine/godot/blob/master/core/math/audio_frame.h) that mixes everything down to a weird 32-bit floating point stereo mix. The mixing might also require going through a compiled language like a GDExtension, as recommended by the [official tutorial](https://docs.godotengine.org/en/stable/classes/class_audiostreamgenerator.html), though a lower mix rate may also be tolerable. The simulation may also have to be offloaded for this to be pulled off. Do it all in GDScript until its working enough to tell if it will be fun or not.
2. Automatic wire placement algorithm. I have this crackpot idea for a wire laying algorithm that uses pathfinding to place wires in a sensible way. The idea is that wires move in right-angled S-shapes to get from point A to point B, but they can use some kind of cost-minimization algorithm that _rewards_ them for going parallel to other wires, and _penalizes_ them for crossing wires. Could be fun.

