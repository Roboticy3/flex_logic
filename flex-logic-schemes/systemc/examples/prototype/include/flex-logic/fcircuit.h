#ifndef _FCIRCUIT_H
#define _FCIRCUIT_H

#include <vector>

#include <systemc.h>

/**
 * Gate base class- empty module which all gates can inherit from.
 */
SC_MODULE(sc_gate) {};

/**
 * Container for sc_gate and sc_signal, managing their relationships
 * The goal of the prototype is not to develop this class, so just keep the
 *  operations simple.
 * 
 * This class will add gates in a "stack" between two inputs and one output.
 * 
 * add_gate(type):
 *   A ----V        =>  A ----V----V
 *   B ---XOR--- C  =>  B ---XOR--type--- C
 * 
 * This is sufficiently complex to test that the simulation is working, because
 *  the serialzed state must be preserved for `type` to have the correct input.
 */
class fcircuit {
  std::vector<sc_gate *> gates;
  std::vector<sc_signal<sc_logic> *> signals;

  //add a gate and connect it to the outputs.
  int add_gate(int type);
};

#endif