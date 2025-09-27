# GHDL
https://deepwiki.com/ghdl/ghdl
Good news, I already have GHDL installed on WSL. My goal is to compile a program, then figure out how feasible it would be to package ghdl in with another software. Here's the guide I was following back then: https://ghdl.github.io/ghdl/quick_start/simulation/index.html

Let's see if I can put the pieces back together.
## Concepts
1. [[Analysis]]: the conversion of "design units" (in this case, VHDL source) to internal representation. Tokenization and Syntax
2. [[Elaboration]]: the conversion of the internal representation to machine code.
3. Run: executing a design to test behavior, generate waveforms.
4. [[VHDL Standard]]: ghdl allows for the selection of vhdl standard (doesn't that technically make it mixed-language??)
5. [[Synthesizable]]: GHDL allows for the construction of non-synthesizable circuits, which goes against the mission of the game. So I've added that to the soft reqs. It should be pretty easy to avoid.

## Evaluation
1. In the heartbeat example, generating and exploring waveforms is all static. It would take work to do that live. In the worst case scenario, I'd have to fork or add GHDL as a module, then access the code somehow to look at the waveform. Go ahead an look if that's possible. Don't dive too deep though. 5 minutes.
	1. I looked at www.deepwiki.com for the ghdl repo. Conversation link: https://deepwiki.com/search/what-would-i-have-to-do-to-int_3369ad5c-7842-4d16-84fa-1bc5abf33277 
	2. The big problem with using this library is that it's written primarily in Ada, which I have no experience in. Once the simulation is done, it may be efficient enough to access it through the python bindings, but it would take a while to get a prototype up.
	3. It appears ghdl is not delta-simulated but exact simulated, which may require some adaptation to meet the requirements.
	4. Input and output would necessarily be in different formats. To work better with the requirements' notion of streams, a converter from some output to the input text streams could be used.
	5. In theory, this would still be very good for debugging. I'd have the functions necessary to build my own checkpointing system. I'll still take this as the best candidate if the other GPL option isn't much better.
2. The "Circuit" runner would have to use the testbench convention. This might actually be very good for the "health" of the codebase.
3. Compiling small circuits is very much instantaneous on modern systems.
# NVC
https://deepwiki.com/nickg/nvc
This looks a lot better than GHDL already. It's written in C, which will be slightly easier for me to work with. I actually recognize some of the fields in the runtime structure, too. However, GHDL also has a python interface that is absent here, and a longer lifespan of community support. Also, my C projects have a habit of going haywire.

What I can do with NVC that I can't with HDL is make a proposed process for figuring out how interface with it. [[Proposed NVC Process]]