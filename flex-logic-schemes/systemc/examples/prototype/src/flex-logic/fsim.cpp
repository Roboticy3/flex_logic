#include <systemc.h>
#include <flex-logic/fsim.h>

void ftestbench::process() {
  size_t step = 0;
  while (1) {
    A.write(input_pattern[step % input_pattern.size()].first);
    B.write(input_pattern[step % input_pattern.size()].second);

    wait();
    step++;
  }
}