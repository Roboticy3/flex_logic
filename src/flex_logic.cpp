#include <queue>

#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>
#include <godot_cpp/templates/vector.hpp>

#include <godot_cpp/core/class_db.hpp>
#include <godot_cpp/core/print_string.hpp>

#include "flex_logic.h"

using namespace godot;

size_t FlexLogic::get_net_count() const {
  return nets.size();
}

/**
 * Table describing the driving of states. Columns represent, in order, each
 * state of the driver. Rows represent each state of the target net. If the
 * order of FlexNetState changes or grows, this must also change or grow.
 * 
 * Table from Buckwell 4-2
 */
const WireState Buckwell42[WireState::MAX * WireState::MAX] = {
//    V0  V1  X   Z   U <- driver state
/*V0*/V0, X,  X,  V0, U,
/*V1*/X,  V1, X,  V1, U,
/*X*/ X,  X,  X,  X,  U,
/*Z*/ V0, V1, X,  Z,  U,
/*U*/ V0, V1, X,  Z,  U,
//^ target state
};

/**
 * Merges the state of `p_from` into `p_to` via the buckwell table.
 * The operation is performed "bitwise", for each element of states.
 */
void FlexLogic::drive(const FlexNetState &p_from, FlexNetState &p_to) {  
  for (int i = 0; i < WIDTH; ++i) {
    p_to.states[i] = Buckwell42[p_to.states[i] * WireState::MAX + p_from.states[i]];
  }
}

void FlexLogic::solver_wire(const Vector<FlexNetState *> &p_states, std::queue<size_t> &r_event_queue) {
  if (p_states.size() < 2) {
    return;
  }
  
  //Copy our value into all connections, then add them to the event queue
  for (int i = 1; i < p_states.size(); ++i) {
    drive(*p_states[0], *p_states[i]);
    r_event_queue.push(i);
  } 
}

void FlexLogic::step_from(size_t p_start_index) {
  if (p_start_index >= get_net_count()) {
    return;
  }

  event_queue.push(p_start_index);
  FlexNetState *state = &net_states[p_start_index];
  
  solvers[state->solver](connections[p_start_index], event_queue);

}

void FlexLogic::_bind_methods() {
  
}

FlexLogic::FlexLogic() {
  
}