Attempts to implement [[Functional Requirements]].

Questions to answer: what HDL environments are out there and what would it take to integrate them with another software? How is it debugged?

The final result will likely be a mix of both approaches. Deepwiki code analysis has proven valuable for looking into the data structures behind different simulation libraries. It may be possible to frankenstein library code and custom code together. I thought it would be more strongly in favor of one or the other, but since it isn't, this honestly isn't in scope for requirements engineering. Still, it was a good use of time and research to dig into the existing tools.
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
