#include <systemc.h>

// Simple SystemC module
SC_MODULE(SimpleModule) {
    sc_in<bool> clock; // Clock input

    void process() {
        std::cout << "Hello, SystemC! Time: " << sc_time_stamp() << std::endl;
    }

    SC_CTOR(SimpleModule) {
        SC_METHOD(process);
        sensitive << clock.pos(); // Trigger on positive edge of clock
    }
};

int sc_main(int argc, char* argv[]) {
    sc_clock clock("clock", 1, SC_NS); // Create a clock signal with 1 ns period

    SimpleModule module("simple_module");
    module.clock(clock);

    sc_start(10, SC_NS); // Run the simulation for 10 ns

    return 0;
}