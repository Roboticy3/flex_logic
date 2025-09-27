From [[Try out HDL Simulators]], here is my idea for possibly integrating NVC:
1. Do NVC tutorials
2. Compile NVC from source with .gdb on
3. Try to inject and extract signals from the runtime using https://deepwiki.com/nickg/nvc. 
4. Try to save and load a state
5. Try to save and load a state onto a modified circuit. 
	1. Since the architecture gets all crazy during compilation, it's probably better to bake the values into the library's construct of the circuit and project them back in from there.
This will also work to test GHDL, since gdb is compatible with Ada.