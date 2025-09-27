#include <vector>
#include <queue>

#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>
#include <godot_cpp/templates/vector.hpp>

#include <godot_cpp/core/class_db.hpp>
#include <godot_cpp/core/print_string.hpp>

#include "flex_state.h"
#include "flex_net.h"
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

void FlexLogic::solver_wire(const std::vector<FlexNetState *> &p_states, std::queue<size_t> &r_event_queue) {
  if (p_states.size() < 2) {
    return;
  }
  
  //Copy our value into all connections, then add them to the event queue
  for (int i = 1; i < p_states.size(); ++i) {
    drive(*p_states[0], *p_states[i]);
    r_event_queue.push(i);
  } 
}

void FlexLogic::restore_connections_from(size_t p_id) {
  net_states[p_id].connections.clear();


  const FlexNet *net = nets[p_id];
  for (const FlexNet *conn : net->get_connections_raw()) {
    size_t i = nets.find(conn);
    if (i != -1) {
      net_states[p_id].connections.push_back(&net_states[i]);
    }
  }
}

void FlexLogic::restore_connections() {
  for (size_t i = 0; i < nets.size(); ++i) {
    restore_connections_from(i);
  }
}

bool FlexLogic::add_net(const FlexNet *p_net, PackedInt32Array p_initial_state) {
  size_t id = nets.find(p_net);
  if (id != -1) {
    return false;
  }
  id = nets.size();

  net_states.push_back(FlexNetState { {WireState::U}, p_net->get_solver(), {}});
  nets.push_back(p_net);

  for (int i = 0; i < WIDTH && i < p_initial_state.size(); ++i) {
    int val = p_initial_state[i];
    if (val < 0 || val >= WireState::MAX) {
      ERR_PRINT("FlexLogic::add_net: Invalid initial state value " + itos(val) + " at index " + itos(i));
      net_states[id].states[i] = WireState::U;
      continue;
    }
    net_states[id].states[i] = (WireState)val;
  }

  return true;
}

bool FlexLogic::remove_net(const FlexNet *p_net) {
  size_t id = nets.find(p_net);
  if (id == -1) {
    return false;
  }

  nets.remove_at(id);
  net_states.erase(net_states.begin() + id);

  restore_connections();

  return true;
}

bool FlexLogic::step_from(const FlexNet *p_start_net) {
  size_t id = nets.find(p_start_net);
  if (id == -1) {
    return false;
  }

  event_queue.push(id);
  FlexNetState *state = &net_states[id];
  
  solvers[state->solver](state->connections, event_queue);

  std::queue<size_t> dump(event_queue);
  while (!dump.empty()) {
    size_t next_id = dump.front();
    dump.pop();
    FlexNetState *next_state = &net_states[next_id];
    solvers[next_state->solver](next_state->connections, event_queue);
  }

  return true;
}

Array FlexLogic::get_nets() const {
  Array arr;
  for (const FlexNet *net : nets) {
    arr.push_back(net);
  }
  return arr;
}

TypedArray<PackedInt32Array> FlexLogic::get_state() const {
  Array arr;
  for (const FlexNetState &state : net_states) {
    PackedInt32Array state_arr;
    for (int i = 0; i < WIDTH; ++i) {
      state_arr.push_back((int)state.states[i]);
    }
    arr.push_back(state_arr);
  }
  return arr;
}

void FlexLogic::_bind_methods() {
  ClassDB::bind_method(D_METHOD("add_net", "net"), &FlexLogic::add_net);
  ClassDB::bind_method(D_METHOD("remove_net", "net"), &FlexLogic::remove_net);
  ClassDB::bind_method(D_METHOD("step_from", "start_index"), &FlexLogic::step_from);
  ClassDB::bind_method(D_METHOD("get_nets"), &FlexLogic::get_nets);
  ClassDB::bind_method(D_METHOD("get_state"), &FlexLogic::get_state);

  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "nets", PROPERTY_HINT_NONE, "", PROPERTY_USAGE_READ_ONLY), "", "get_nets");
  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "state", PROPERTY_HINT_ARRAY_TYPE, "PackedInt32Array", PROPERTY_USAGE_READ_ONLY), "", "get_state");

}

FlexLogic::FlexLogic() {
  
}