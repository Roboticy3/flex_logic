Issues with old Lageta:
1. Bland
2. Technologically limited by the platform

Solutions:
1. Bash it with Pikuniku, also make it a.. music game!
2. Make it in Godot

# Simulations, Levels, and Layers of Abstraction
Levels are organized into worlds. The player can move smoothly between levels, sub-levels, worlds, and more by zooming in and out. I've already implemented this, actually. Though it's janky. The point of this mechanic is to emphasize how computing is done in layers.
## Levels
Levels can contain circuits with inputs and outputs. If they do, they can be simulated with different inputs both inside the level, and in the next world up, where the inputs and outputs of levels will also be accessible.

Levels can also have clear conditions, desired strings out outputs from the simulation. Once reached, the level is cleared, unlocking more levels, gates, and tools.
## Circuits
Circuits have inputs, outputs, gates, and wires. They are simulated continuously from whenever the inputs change, to whenever they converge. Each tick corresponds to a short amount of "real" time, emulating real gate delay.

Circuits come with a signal reader. The signal shows a timeline of inputs, desired outputs, and a play button. When the player presses the play button, the timeline progresses, sending inputs to the circuit in order and overlaying the actual outputs over the desired outputs. Signal reader channels, for the sake of tutorialization, can also float around and position themselves near the actual pins that they represent.

Gates will be displayed in an inventory bar, organized into folders that open on hover. Gates can be invoked by clicking or with hotkeys and placed on the circuit by clicking. Wires can be drawn between them by clicking and dragging.

Another visual aid will be a power and ground line at the top and bottom of the circuit that can be hidden with a hotkey. They will draw vertical wires through any gates that are placed in the circuit, giving some diegetic reason to how they are receiving power, and will glow when signal is being played.
## Tools
### Probe
A metal probe that can be placed on any wire or pin of the circuit to display an extra output channel. Unlocked in world 1.
### Powered probe
An upgraded version of the probe that can power pins and shock creatures. Unlocked whenever Transistors are introduced.
### Bus
A "fake" tool, in that it has no implementation of its own but won't have any uses until the right gates are unlocked. Some gates will convert between *n*-pin inputs/outputs and "bus" pins. Bus pins will carry bundles of wires as a single wire. Technically, everything should be a "bus" of width 1 or more. The underlying simulation would be easiest to implement this way, with wires and pins storing an integer or big integer for their state, and a bit-width determining what they should be allowed to connect with. Unlocked at the end of world 3

Wide wires will display their value as an integer embedded in the center of the wire.
### Big Hands
Allows the user to manipulate and draw wires on the level select screen. Not just inside levels themselves. Unlocked at the start of world 3.
### Bus Probes
A special probe, unlocked after learning decoding, that can decode bus wires and pins into integers, strings, or hex. Probably won't be in the demo.
### Multiselect
QoL. Allows for clicking and dragging to select multiple gates, as well as ctrl-clicking to copy the current selection, or just the gate under the mouse. Probably won't be in the demo.
# Simulations, but awesome
I have five worlds planned, covering the same content of the original. I already have the world select screen mostly set up, with the exception of the last world.
## World 1: Tutorial
The first will only be a few levels. It will be the first screen the player sees and contain tutorial info on the user interface. After completing this world, worlds 2 and 3 will become available.

This is also where the main hook of the simulation will be introduced. Probes and outputs will be audible as sound! At first, only 1-bit square waves, likely scaled up to cooperate with audio drivers, will play. But different worlds and levels will have different audio effects applied to them, and the sound will become more integrated with the gameplay as the player progresses. 

In this demo, the signal viewer should be set up so that the *volume* of the audio appears as signal, and is used to process the circuit to. In the long run, I want audio to be passed through the circuit bitwise, and I want more advanced levels to really dig into that. But I'm not sure how feasible that will be performance-wise. Godot does have primitives for generating audio streams on the fly, though. It may even be easier to take this approach from the start. Ideally, circuits will be implemented oblivious to how the audio will interact with it, then the audio will be implemented on top of them and the circuits will be tuned to account.
## World 2: Combination
Next, there is a level exploring all the wonderful things one can do with And and Or gates. Most of these levels are optional. One level will unlock the Xor gate, and technically* be available from the start. The player can complete this immediately to unlock World 4. Or they can complete the levels in order, progressively unlocking more and more combinations until reaching Xor. more naturally.

So, I said technically. Levels can simulate the circuits that are inside of them. Levels are only "locked" because they do not have the right arrangement of wires connected to them. If they are powered in a specific way, a way which is achieved automatically by simply filling in the circuits of previous levels, they become available. But, the user can shortcut this using the Big Hands tool.

If the user completes Xor level without having completed World 3, they will be shown a cutscene highlighting the last level of World 3.
## World 3: Coding
World 3 is all about encoding and decoding signals. This is where the user will learn about binary numbers, and busses, and gain the Big Hands tool needed to skip ahead in World 2 Completing world 3 unlocks world 5.

In introducing busses. This game also introduces higher fidelity, mouth-watering *8-bit* audio tracks! With triangle waves and sine waves being built from square-wave bricks. However, in a volume-only audio simulation, only stacks of square waves, with increasing and decreasing volume, would be possible. There could also be one optional level for unlocking 16-bit splitters/joiners with 16-bit audio.
## World 4: Arithmetic
The big moment of this demo is going to be an 8-bit adder that can add 2 audio tracks together to "mix" them. The arithmetic world will introduce carries and build up to a concert at the end, with additional levels for 3-track mixing and 16-bit mixing.
## World 5: Decimal
I'm not sure if I want to include this, because it isn't very musical. The player will build up to a binary-to-7-segment encoder, testing their mastery of busses and coding. Maybe the bus probe would be a good reward for this? It's super bonus.
# Transistors running amuck
Picture this: the player gets up to drink water, only to come back to find their levels have been locked! They check the one of the remaining available levels, only to find that their circuit has been tampered with behind their backs! 

After unlocking the Big Hands tool, Transistors, creatures with two legs, a circular body, and one arm out the top, will start to occasionally appear. Sometimes the player will catch them tampering with the level they are currently working on.

Transistors walk across the bottom of the screen, then tower on top of each other, forming two walls. When each wall reaches a pin, they'll bite, cutting both ends of a wire and taking it with them for dinner! The wire will become a rope strung between the two biters, and they'll run off-screen.

You can try to interrupt them by placing gates in the way of the walls, or by running current through wires to fry them. If you place your wires close to the bottom of the circuit, they should be able to take them in as little as 5 seconds. 

For every session in which you attack or hurt them, your reputation decreases. The likelyhood that they will attack levels in the background (without any of the physics sim, lol. The wires just disappear), will increase. This will be indicated by a "repuation lost" message flashing up in the center of the screen when they are harmed.

To regain reputation ("reputation gained" message), simply allow them to eat wires and replace the wires as needed. This will force the player to remember where they put their wires, often more than one at a time, in levels they completed a while back. 

Sometimes, they will appear on world screens as wec. Again, walking in from the sides along the bottom of the screen. But they are not aggressive here. Instead, if reputation is positive, they will try to reach for the user's cursor and give you a high five. Otherwise, they will run from the cursor. They can still be shocked by current from the level select screen.

Finally, with positive reputation, in the World 3 levels with 8 bit music, and the adder/layering levels, they will start a dance party when the player inputs the correct solution. If the solution has succeeded after a couple seconds of signal, Transistors will begin to run in and starting gettin down. If the signal "breaks" (a case comes up where the user's adder returns the wrong result), a record-scratch will play and the Transistors will hang their heads, collapse out of their dance-walls, and go home. This will also result in a point removed from reputation. This links them back to the puzzle-solving loop by making the puzzle solution have an effect on their behavior.

It might be good to have individual reputations, as wec. Just temporary ones. So you could, for example, pick a Transistor up and have it scurry away when you let go.
# Goblins, and boomboxes.
Another creature inhabiting this world is the relatively rare Goblin. The player does not maintain a reputation with Goblins. It may be good to script their appearances, since they won't have a "motivation" to attack the player otherwise.

First, the player will hear a trap beat playing quietly in the background. Transistors may start to fill in in front of it as it gets louder, before the Goblin walks in slowly with the beatbox over its shoulder. At this point, their music will drown out anything playing out of the circuit, and stop it from playing.

The boomboxes attract transistors wherever they go, attacking your wires if they don't like you. But, worst of all, goblins eat gates whole! Walking across the bottom of the screen, they look up.. Their arm extends vertically, then turns around to a pressing finger, while their jaw expands, tracing out a perfect gate-shaped rectangle to slowly push the target into. Shocking them will deter them eventually, but not immediately. If shocked while eating a gate, the gate will stay displaced, and they can only eat one gate at a time.

While not eating, they'll bump to the music playing out of their boombox, eventually putting the boombox down to dance a little, for about 2 seconds out of every 10 seconds. If the player clicks on them while carrying the boombox, they'll hunch over it. But if the player grabs the box while its on the ground, they're free to carry it wherever. It's well known that goblins cannot jump, so perch the boombox on a gate to keep it out of reach, or throw it off screen. Once it has fallen off screen, it will fade out, any Transistors will follow it away, and the Goblin will hang its head and go home. The circuit will become playable again too.