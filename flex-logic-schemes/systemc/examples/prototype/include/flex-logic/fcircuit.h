#ifndef _FCIRCUIT_H
#define _FCIRCUIT_H

#include <vector>

#include <systemc.h>

/**
 * Gate base class- empty module which all gates can inherit from.
 */
SC_MODULE(sc_gate) {};

class sc_xor : public sc_gate {
  public:
    sc_signal<sc_logic> A;
    sc_signal<sc_logic> B;
    sc_signal<sc_logic> C;

    void process() {
      C = A ^ B;
    }

    SC_CTOR(sc_xor) {
      SC_THREAD(process);
      sensitive << A << B;
    }
};

class sc_and : public sc_gate {
  public:
    sc_signal<sc_logic> A;
    sc_signal<sc_logic> B;
    sc_signal<sc_logic> C;

    void process() {
      C = A & B;
    }

    SC_CTOR(sc_and) {
      SC_THREAD(process);
      sensitive << A << B;
    }
};

class sc_or : public sc_gate {
  public:
    sc_signal<sc_logic> A;
    sc_signal<sc_logic> B;
    sc_signal<sc_logic> C;

    void process() {
      C = A | B;
    }

    SC_CTOR(sc_or) {
      SC_THREAD(process);
      sensitive << A << B;
    }
};

enum sc_gtype {
  XOR,
  AND,
  OR
};

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

  public:
    //add a gate and connect it to the outputs.
    int add_gate(sc_gtype type);

    fcircuit();
    ~fcircuit();
};

#endif