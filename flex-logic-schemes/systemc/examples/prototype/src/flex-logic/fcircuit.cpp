#include <string>

#include <systemc.h>

#include <flex-logic/fcircuit.h>
#include <flex-logic/fsim.h>

fcircuit::fcircuit() {

}

int fcircuit::add_gate(sc_gtype type) {
  //add the gate to the circuit
  sc_gate *gate;
  sc_module_name label = std::to_string(type).c_str();

  switch (type) {
    case XOR: gate = new sc_xor(label);
    case AND: gate = new sc_and(label);
    case OR: gate = new sc_or(label); 
  }

  //TODO: connect the signals

  gates.push_back(gate);

  return 0;
}

fcircuit::~fcircuit() {
  for (auto g : gates) {
    delete g;
  }

  for (auto s: signals) {
    delete s;
  }
}