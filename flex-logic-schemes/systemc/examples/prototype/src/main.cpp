#include <iostream>

#include <assert.h>
#include <systemc.h>

#include <flex-logic/fcircuit.h>
#include <flex-logic/fsim.h>

/**
 * Tests:
 * 
 * Simulation behaves as expected with no modifications
 * 
 * Simulation behaves as expected after modification and propogation
 * 
 * Simulation behaves as expected after modification and resuming
 * 
 * Benchmarking for large circuits.
 */

void no_modification_xor() {
  fcircuit circuit;
  fsim_proto sim(circuit);
  
}

void no_modification_and() {

}

void no_modification_or() {
  
}

int sc_main(int argc, char* argv[]) {
  std::cout << "Hello World\n";

  no_modification_xor();
  no_modification_and();
  no_modification_or();
}