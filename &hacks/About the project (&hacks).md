## Inspiration
2 years ago, I started making Lageta, a puzzle game to test the effectiveness of educational games in undergraduate digital architecture. Under the mentorship of Professor Yifan Sun, we completed the project in about 9 months. However, it had a few issues:
 - Accuracy of digital logic simulation
 - Poor visualization

It's really hard to visualize digital logic simulation in an accessible way. I had an idea later to use the Godot game engine's dynamic audio feature to instead use user-created circuits to generate _audio_. However, that would require a digital logic simulation that could run at at least 11kHz for small circuits, which just isn't possible with the propagation-based JavaScript system that Lageta runs on.  So, I started work on a Godot C++ module to get the job done.

Now I'm at an issue of complexity. After reading through relevant parts of the Verilog standard, researching existing simulators, and writing a lot about how I want an audio-focused digital logic game to play, I realized I need a solid design spec to do what I want to do.

Everything on the market is just a couple features shy of what I need. Since I'm making a game, the circuits have to be editable in real-time in some form. Since I'm using audio feedback, the circuits also have to use low-level abstractions to run efficiently. These two specifically do not typically align with each other. Usually a simulator is either compiled or dynamic. The final implementation will likely use a live-recompilation system.

Recently, I started taking Software Engineering as an elective. It talks about the merits of requirements engineering and technical specification. I thought it would be good to try and follow the lessons from that class, so I stayed up late last week writing a list of User Stories for the core library.

## What it does
I'm making a _technical specification_ for flex-logic, a library with Godot integration for simulating digital logic circuits. It would take months to actually code the thing, but since I already have the user stories, a technical, interface level specification might be in reach.

## How we built it
Obsidian :)

## Challenges we ran into
Digital logic circuits are notoriously hard data structures to pin down, even under slight variations to their requirements.

## Accomplishments that we're proud of
The concise set of use cases produced for this project provide a great outline for dynamic circuit simulators everywhere, even those outside the reaches of digital logic.

## What we learned
Requirements engineering and specification 

## What's next for Flex Logic Spec
Keep working on it. Hopefully get the simulator done. The end goal for the spec itself is to integrate it into a larger game design document going into how the simulator can be used to make a compelling experience.
