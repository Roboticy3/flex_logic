Nets are the data structure for carrying circuit events in any self-respecting simulator. They've been pretty tough to get my head around. I'm putting a lot of work into implementing them, and I think they should play a role in the game.

There is also a problem where the concept of wires and gates are hard for new players to wrap their heads around.

What if the first few levels the player is exposed to have no gates or wires at all?

Instead, they control a squishy blob on a grid of pins. The blob corresponds to a net that moves one pin in any direction when the player presses WASD. This will make the player identify with the net, and be more aware of their position. 

The first level is then simply moving the net into a position that connects an input and an output. A good allegory is how, in *Portal*, the portals are static before giving the player the portal gun. By making the gates static for a while, we're showing the player how gates work before giving them the ability to place gates.

The first level has an input gate brightening and dimming in response to an audio signal, which starts playing when the net moves into place. The connection starting could even be accompanied by an artificially added amp-plug/unplug sfx. Two more levels are unlocked.

Probes are achieved by simply moving over the pin that you want to hear. 

The second level has an AND gate and two inputs. The first input is music. The second is an LFO signal. Under probe, LFO signals make a pure square tone when they are high, and a silent (or low? whatever sounds better) tone when low. Either signal produces noise in an illegal state. When the player connects both inputs, they are prompted to press Space to *drop* the net. If they do this, the net is dropped and another net spawns in, which they can use to connect the AND gate to the output.

If one net walks into a dropped one, they merge, and the player gains control of the whole mass.

 > Careful. Players will find this confusing. I think players won't know what's going on when they beat the second level. There has to be something else to indicate what they are doing.

Maybe, to make things easier, give the player a button to hear the desired signal? It's clear, but it feels lame. Keep an eye on this.

Here's another idea: the third level will have the player control an AND gate. The level contains a similar layout to the completed second level, but and OR gate is where the AND gate should be, and the net connecting it to the output is twice as tall. The player must push the OR gate out of the way using the AND gate. For one step, both gates will be wired to the net, causing a conflicting state, which will sound like white noise. The goal is to introduce that the type of the gate can be important, while also exposing the player to differently shaped nets and conflicting states. The player can also push the OR gate incorrectly, and have to undo their move, which they will be prompted for if they push the OR gate the wrong way.

After the third level, the player unlocks truth tables. Every gate they see will come with a truth table describing its behavior under different inputs, and they can hover over a gate to see its table. The table will show a picture of the gate for each combination of high/low inputs. In reality, the inputs are 16 bit unsigned. These will display in-game in a gradient based on what value they decode to. The gradient should be the same as the decibel gradient used in a standard DAW. Nets will also take the color of their state.

The Goblin mechanic from [[Lageta 2]] ports nicely to this format. Transistors might be able to turn nets into softbodies instead of ropes.. Idk, that's a stretch.
# Questions
1. Can you pull this off?
	1. I actually think this will be easier to pull off than my current trajectory. I don't have to implement a wire controller. I should be implementing and testing the simulation, currently. This would motivate me more in that direction. I would need:
		1. Simulation controller
		2. Sokoban controller
		3. Cooperation between the components
		4. Basic graphics and audio demonstrating functionality in godot
2. How will you reconcile a more abstract puzzle tone with the established setting and story? 
	1. In this iteration, the nets controlled by the player are creatures deployed from their ship. The player is of the same species, like a commander and their soldiers. 
3. What is the end goal?
	1. I want to make a few more levels after those described, ending off with the ability to place gates as well as nets. But I might end up going full sokoban on this.