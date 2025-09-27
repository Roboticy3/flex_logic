Attempts to implement [[Functional Requirements]].

Questions to answer: what HDL environments are out there and what would it take to integrate them with another software? How is it debugged?

# 1: Dynamic Recompilation
Main idea:
 - When a change is made to the circuit, compile to VHDL or another HDL and run.
 - Primary goals are to simplify and standardize the codebase. If the weaknesses start to outweigh this by making the code too complex, other options should be sought out.

| Strengths                                                                                                                                            | Weaknesses                                                                                                                                                                                                     |
| ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| SPEED                                                                                                                                                | Changing while running is hard. Depending on the VHDL interpreter in use, it may require recompiling, then rerunning the circuit up to the pause point, which would also be extra hard with substep debugging. |
| Standardized and realistic simulation behavior.                                                                                                      | Debugging is hard. More research required to understand available VHDL debugging options.                                                                                                                      |
| Standardized circuit description. With annotation, circuits could even be loaded from their VHDL descriptions. Though that might not be a good idea. |                                                                                                                                                                                                                |
Fully cutting the pausing behavior is not off the table if this turns out the be a viable solution.
# 2: Custom Simulator
Main idea:
 - Get more control over the simulator's implementation with a custom simulator
 - Primary goals are more control over the simulator, which will make it easier to make design tweaks as needed when the simulator doesn't "feel right" in game.

| Strengths                                                                   | Weaknesses                                                                                     |
| --------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| Possibly better portability than a simulator that relies on external tools. | Hard to implement a realistic simulator.                                                       |
| All debugging tools for regular programming are at my disposal.             | Hard to match standards of other simulators. I'm not enough of a field expert to trust my gut. |
| Easier to adjust the simulator for design purposes.                         |                                                                                                |
# 3. Open Source Dynamic Simulator
Main idea:
 - Mooch off of open source. Basically CircuitVerse
 - Primary goal is to offload the implementation of the simulator without relying on a compiled system


| Strengths                                                        | Weaknesses                                                            |
| ---------------------------------------------------------------- | --------------------------------------------------------------------- |
| Wide array of existing features that can be directly integrated. | Current options have limited performance compare to compiled systems. |
