## Inspiration
2 years ago, I started making [Lageta](lageta.netlify.app), a puzzle game to test the effectiveness of educational games in undergraduate digital architecture. Under the mentorship of Professor Yifan Sun, we completed the project in about 9 months. However, it had a few issues:
 - Accuracy and speed of digital logic simulation
 - Poor visualization

It's really hard to visualize digital logic simulation in an accessible way. I had an idea later to use the Godot game engine's dynamic audio feature to instead use circuits to generate _audio_. However, that would require a digital logic simulation that could clock at least 11kHz for small circuits, which just isn't possible with the TypeScript event propagation system that Lageta runs on.  So, I started work on a Godot C++ module to get the job done.

I quickly ran into an issue of complexity. It was hard to design a workflow where I could quickly test and modify my code because of how Godot plugins are architected. I was also doing a lot of research into digital logic simulation and quickly finding that it was more than I could handle in a few weekends. 

I was also looking at existing solutions, but most systems on the market are a couple features shy of what I need. Since I'm making a game, the circuits have to be editable in real-time in some form. Since I'm using audio feedback, the circuits also need to use low-level abstractions to run efficiently. These two features specifically do not typically align with each other. Usually a simulator is either compiled or dynamic. Not both.

Recently, I started taking Software Engineering in school as an elective. It talks about the merits of requirements engineering and technical specification. I think these techniques might have what it takes to defeat the problem, at least on paper. So, I stayed up late last week writing a list of User Stories for the core library.
## What it does
This weekend, I'm making a _technical specification_ for flex-logic, a library for simulating digital logic circuits. It would take months to actually code the thing, but since I already have the user stories, a technical specification at the interface level is in reach. Hopefully that's what I can present to you by noon on Sunday, along with the other design artifacts I made in the process.

All artifacts from this weekend are in the `&hacks` folder of the [github repo](https://github.com/Roboticy3/flex_logic). You can also view details of the project in the [github project](https://github.com/users/Roboticy3/projects/2).
## How we built it
Obsidian :)
## Challenges we ran into
Digital logic circuits are notoriously hard data structures to pin down, even under slight variations to their requirements. Difficult research made up the bulk of the work, including complex theory and esoteric source code. The hardest design questions to tackle were:
1. Who is the library for? Is it just for me, or is the library meant to be portable?
2. Should I design a custom-made simulator or a system that utilizes existing simulation tools like GHDL or NVC?
3. Is it worth it to go for an interface level specification?
## Accomplishments that we're proud of
The concise set of use cases, generated from the user stories, provided a great outline for dynamic circuit simulators everywhere, even those outside the reaches of digital logic.
## What we learned
Requirements engineering and specification management are great ways to keep a high level view of a project and to hone your design skills without getting lost in the weeds.
## What's next for Flex Logic Spec
Keep working on it. Hopefully get the simulator done concretely. The end goal for the spec itself is to integrate it into a larger game design document going into depth on how the simulator can be used to make a compelling experience. You can read the original pitch doc for in `Lageta TWO/Technical Breakdown`.
