In the [[Lageta 2|Original Design Document]], I outline some creatures to inhibit the player. Here are some other thoughts I've had about them.

# Light ecosystem mechanics
In each world, at a given time, there can be a total count of each creature type, influencing how often they appear. These populations can then be regulated over time using dynamics. 

Creatures I want in the game:
 - Transistors
 - Goblins
 - Cats
 - Ancients
 - Hain

Cats eat Transistors.
Ancients eat Cats and Goblins, but do not reproduce.
Hain randomly increase and decrease in population, and "scare" all other creatures but does not predate them.

If a creature is currently inside the active level, and becomes scared, they will run away. Now, lets talk about these new creatures. All creatures become scared when a predator creature appears.

Creatures also stay in a level when the player leaves, and can enter levels with limited simulation abilities without the presence of the player. Levels with creatures in them will show icons describing which creatures are present and how many of each in the top left of their card on the level select screen.
# Cats
Cats are harmless to the player, and jump onto different components. They have different responses to stimuli, but these responses are largely aesthetic. Unlike real life, they are immune to electrocution, and instead seem to enjoy a steady current, glowing yellow and smiling. If a cat sees a singular transistor, it will play a pounce animation and eat it. However, it is scared of stacks of transistors. If no stimuli comes for 1 minute, and it hasn't wandered off, the cat will find a place to sit or fall asleep. Resting cats will only move if another stimuli happens within a very short range.
# Ancients
Ancients are blue, tall, humanoid creatures capable of speech. They will not eat in front of the player, nor explain what they eat, but those creatures are still scared of them. They don't do anything visible except control population, and walk across the screen, usually without stopping. They dislike goblin boomboxes, and will stop to sneer at them before scaring them off. It should be very rare for more than one to appear at a time.

Some Ancients will be scripted to talk to the player- appearing in custom scenes in-between levels similar to tutorials. Dialogue may have choices.

Ancients all appear male, but have different tones of blue and body proportions. They dislike transistors, but do not eat them, and transistors are not scared of them. They may try to chase transistors off, even confronting towers, but the transistors do not fully run away, only scattering.

It should be a twist that they eat goblins and cats. Initially, they may appear virtuous for helping the player with puzzles, and having a "noble" demeanor with their good posture and relatively high intelligence. When it is revealed that they eat these other creatures, it should make the player question their authority. In dialogue, they will be named, and some should be characterized as only eating one or the other. 
# Hain
Hain are information-eating monsters inspired by the modern depiction of the Selk'nam in South America. The story goes that a German priest asked the native Selk'nam of southern Patagonia to show him their rituals, and the men basically made something up that was roughly inspired by a tradition that hadn't been practiced in the area for generations. The result was this sort of "artificial" culture. The supposedly genuine practice was so captivating that it instantly became part of the South American identity. Pretty much anywhere in south South America, you can buy cartoon depictions of these body-painted men as fridge magnets or sculptures. The idea of Hain flips all that on its head.

Hain destroy information by finding real things and making them fictional. They have a slightly lovecraftian feel and should not be introduced until later in the game. They are preceded by a tutorial from an Ancient, describing them.

Hain never appear out of nowhere. First, invalid utf characters (referred to by the Ancients as "entropy particles") start to flow in from one side of the screen for 30 seconds, building up to a crescendo. Then, the creature appears. After this point, the player has 1 minute to beat the level. If they leave, the timer keeps ticking! If this requirement is not met, the level goes black, and cannot be entered for 10 seconds. After this time, the player permanently loses one gate type at random, and the save data for the level which unlocked it. Essentially, the level is un-completed.

This type can be recovered by returning to the associated level and re-completing it. Or the player can continue without recovering the gate.

If a Hain appears with no player present, it simply waits in the level until a player appears for 5 minutes. If the player enters during this period, the 1 minute death timer begins. 

I also saw this shop in Santiago called poloduranceramista that has these awesome clay sculptures with a similar eeriness to them. Like animal-human hybrids made only from primitive shapes and curves.