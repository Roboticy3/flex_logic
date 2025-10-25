#include <systemc.h>

/**
 * Example of using systemc to build a circuit programmatically
 * 
 * Xor gate constructed from nand gates
 * 
 * Oh, you're wondering why all the ports are integers? 
 * 
 * Yeah.. Uh. idk. Full simulation would probably have a custom "logic9" type to
 *  better match vhdl.
 */

SC_MODULE(Gate) {};

class Nand : public Gate {
  public:
    sc_in<int> A;
    sc_in<int> B;
    sc_out<int> C;

    void process() {
      C = ~(A & B);
    }

    SC_CTOR(Nand) {
      SC_METHOD(process);
      sensitive << A << B;
    }
};

SC_MODULE(Circuit) {
  //ports
  sc_signal<int> A;
  sc_signal<int> B;
  sc_signal<int> C;

  //clock
  sc_clock clock;

  //stream
  std::vector<std::pair<int,int>> input_pattern = {
    {0,0}, {1,0}, {0,1}, {1,1}
  };

  //gates
  std::vector<Gate *> gates;

  //nets
  std::vector<sc_signal<int> *> signals;

  //simulation loop
  void process() {
    size_t step = 0;
    while (1) {
      std::cout << "stepping..\n";

      sc_time time = sc_time_stamp();

      A.write(input_pattern[step % input_pattern.size()].first);
      B.write(input_pattern[step % input_pattern.size()].second);

      // Print current circuit state. Thanks ChatGPT :)
      std::cout << "Inputs: A=" << A.read()
                << "  B=" << B.read() << std::endl;

      for (size_t i = 0; i < signals.size(); i++) {
        std::cout << "n" << i + 1 << "=" << signals[i]->read() << "  ";
      }

      std::cout << "Output: C=" << C.read() << std::endl;

      wait();
      step++;
    }
  }

  SC_CTOR(Circuit) : clock("clock", 1, SC_NS) {
    SC_CTHREAD(process, clock);
    sensitive << A << B;

    sc_signal<int> *n1 = new sc_signal<int>("n1");
    sc_signal<int> *n2 = new sc_signal<int>("n2");
    sc_signal<int> *n3 = new sc_signal<int>("n3");
    signals.push_back(n1);
    signals.push_back(n2);
    signals.push_back(n3);

    //circuit construction
    Nand *nand1 = new Nand("nand1");
    Nand *nand2 = new Nand("nand2");
    Nand *nand3 = new Nand("nand3");
    Nand *nand4 = new Nand("nand4");
    gates.push_back(nand1);
    gates.push_back(nand2);
    gates.push_back(nand3);
    gates.push_back(nand4);

    // Wiring
    nand1->A(A);
    nand1->B(B);
    nand1->C(*n1);

    nand2->A(A);
    nand2->B(*n1);
    nand2->C(*n2);

    nand3->A(B);
    nand3->B(*n1);
    nand3->C(*n3);

    nand4->A(*n2);
    nand4->B(*n3);
    nand4->C(C);
  }

  ~Circuit() {
    for (Gate *g : gates) {
      delete g;
    }

    for (sc_signal<int> *s : signals) {
      delete s;
    }
  }
};

//program
int sc_main(int argc, char* argv[]) {

  Circuit tb("testbench");


  sc_start(5, SC_NS);

  return 0;
}