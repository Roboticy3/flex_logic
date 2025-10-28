#ifndef _SYSTEMC_H
#define _SYSTEMC_H

#include <vector>

#include <systemc.h>
#include <flex-logic/fcircuit.h>


/**
 * Code stolen from examples/programmatic-circ
 * Constructor code for building ciruit statically has been commented out
 * Goal is to replace that to some degree with dynamic circuit construction
 * Don't get too crazy with those capabilities, maybe just add gate and add conn
 *  for now.
 */
SC_MODULE(ftestbench) {
  //ports
  sc_signal<sc_logic> A;
  sc_signal<sc_logic> B;
  sc_signal<sc_logic> C;

  //clock
  sc_clock clock;

  //stream
  //in this prototype, the stream should loop until further improvements are made
  std::vector<std::pair<sc_dt::sc_logic, sc_dt::sc_logic>> input_pattern = {
    { sc_dt::sc_logic('0'), sc_dt::sc_logic('0') },
    { sc_dt::sc_logic('1'), sc_dt::sc_logic('0') },
    { sc_dt::sc_logic('0'), sc_dt::sc_logic('1') },
    { sc_dt::sc_logic('1'), sc_dt::sc_logic('1') }
  };

  /**
   * interestingly, we no longer need a direct reference to the gates or signals
   * here, all the good stuff is hidden in the systemc lib.
  //gates
  std::vector<sc_gate *> gates;

  //nets
  std::vector<sc_signal<sc_logic> *> signals;
  */
  
  public: //unfortunately, I think this has to be public.
    void process();
  
  SC_CTOR(ftestbench) : clock("clock", 1, SC_NS) {
    SC_CTHREAD(process, clock);
    sensitive << A << B;
  }
};

/**
 * According to the goal prototype first introduced in #16:
 * "fsim pauses cocotb (systemc now) and serialized its state"
 * ...
 * "fsim compiles the circuit to vhdl". 
 * "fsim sends the vhdl to cocotb..."
 * These two steps can be replaced by:
 *  "fsim converts between fcircuit and SC_MODULE"
 * "fsim starts cocotb (systemc now) with the serialized state"
 * "fsim signals [fcircuit] that the edit has been successfully made"
 * 
 * Keep it simple for now: fsim_proto will simulate the circuit automatically.
 *  When fcircuit invalidates, synchronize_tb will run and perform the actions
 *  above, without resuming simulation.
 * This should then be tested in a variety of situations by main.cpp
 * 
 * wait0_tb will wait until the changes resolve (gates will have no delay for
 *  now).
 * 
 * resume_tb will resume the simulation.
 */
class fsim_proto {
  private:
    ftestbench tb;
    fcircuit circuit;
    std::vector<sc_logic> cache;

    int synchronize_tb();

    int wait0_tb();

    int resume_tb();
};

#endif