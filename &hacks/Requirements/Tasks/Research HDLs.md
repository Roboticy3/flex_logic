As asked in [[Possible approaches]], what relevant information do I have about HDLs?

From 
https://en.wikipedia.org/wiki/Hardware_description_language
and
https://en.wikipedia.org/wiki/List_of_HDL_simulators

After the [[Try out HDL Simulators]] task, I recommend GHDL or NVC as the backend simulator. Both also have relatively simple linkers for linking a design with prebuilt gates, which also means user-generated gates could be compiled, saved, and reused in theory.
# Active-HDL Student Edition
This is a free software suite developed by ALDEC. Technically, it includes all the functionality I could want out of my own project and more. But that doesn't mean I can just drop it in. I need to verify that:
1. It can be used as a library. Looking at the system files, it's definitely possible. It's very well organized into dlls that can be referenced on windows, but this could become an issue if the library is moved or needed on another OS. The free version is Windows only.
2. I can intercept the simulator
It would also be nice to check if that's legal. Honestly, I should assume it's not, since these features are also a produce of Riviera-PRO.
# Incisive/NCSim/irun
This is a proprietary linux library for testing logic controllers. Again, it's proprietary, so a drop in use is probably not ok. But it would be good to try it out on WSL to see what the workflow for simulation is like. Add a task for it, though it might be harder to try out than PyMTL 3, and somehow even more out of date. A limitation of using industry tools is that they may all be towards the older side.
# PyMTL 3/Mamba
Open source python simulator reliant on `verilator`.. which is reliant on C++. So at least it should have good speed. Interface is confusing, but still the strongest candidate so far. The main weakness here is the number of levels code has to go through before it can be recompiled. I would also be doing graph2code from scratch, which I didn't even know was something that could be done for me until Active-HDL. Include in the tryout task.
# GHDL and NVC
Here's that good shit. GPL2+, VHDL-2019 support. It doesn't get much better than that. The rest of the options are "mixed-language" simulators, meaning they support multiple HDLs, sometimes in the same design. While that's an incredible feature, it's utterly unnecessary for a project like this, since I'm focusing on creating an open-ish library that's going to be mainly interacted with through a GUI. These open source tools are well documented, up to date, and free to use. Include in the tryout task.