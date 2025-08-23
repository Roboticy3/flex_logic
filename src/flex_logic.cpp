#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>
#include <godot_cpp/templates/vector.hpp>

#include <godot_cpp/core/class_db.hpp>
#include <godot_cpp/core/print_string.hpp>

#include "flex_logic.h"

using namespace godot;

void FlexNet::solver(Vector<FlexNet *> &r_event_queue) {
  //Copy our value into all connections, then add them to the event queue
  for (FlexNet *connection : connections) {
    drive(connection);
  } 

  r_event_queue.append_array(connections);
}


/**
 * Table describing the driving of states. Columns represent, in order, each
 * state of the driver. Rows represent each state of the target net. If the
 * order of FlexNetState changes or grows, this must also change or grow.
 * 
 * Table from Buckwell 4-2
 */
const FlexNetState Buckwell42[FlexNetState::MAX * FlexNetState::MAX] = {
//    V0  V1  X   Z   U <- driver state
/*V0*/V0, X,  X,  V0, U,
/*V1*/X,  V1, X,  V1, U,
/*X*/ X,  X,  X,  X,  U,
/*Z*/ V0, V1, X,  Z,  U,
/*U*/ V0, V1, X,  Z,  U,
//^ target state
};

/**
 * Merges the state of this FlexNet into `r_to` via the buckwell table.
 * The operation is performed "bitwise", for each element of states.
 */
void FlexNet::drive(FlexNet *r_to) {  
  if (!r_to) return;

  for (int i = 0; i < _get_size_internal(); ++i) {
    r_to->states[i] = Buckwell42[r_to->states[i] * FlexNetState::MAX + states[i]];
  }
}

void FlexNet::set_value(int p_value) {
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (p_value >> i)
      states[i] = FlexNetState::V1; 
    else
      states[i] = FlexNetState::V0; 
  }
}

int FlexNet::get_value() const {
  int value = 0;
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    value |= ((states[i] & 0x01) << i);
  }
  return value;
}

void FlexNet::set_u(int p_mask) {
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (p_mask >> i)
      states[i] = FlexNetState::U;
  }
}

int FlexNet::get_u() const {
  int value = 0;
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (states[i] == FlexNetState::U)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_x(int p_mask) {
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (p_mask >> i)
      states[i] = FlexNetState::X;
  }
}

int FlexNet::get_x() const {
  int value = 0;
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (states[i] == FlexNetState::X)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_z(int p_mask) {
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if ((p_mask >> i) & 0x01)
      states[i] = FlexNetState::Z;
  }
}

int FlexNet::get_z() const {
  int value = 0;
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    if (states[i] == FlexNetState::Z)
      value |= (1 << i);
  }
  return value;
}

size_t FlexNet::get_size() const {
  return _get_size_internal();
}

bool FlexNet::add_connection(FlexNet *p_connection) {
  if (p_connection && !connections.has(p_connection)) {
    connections.push_back(p_connection);
    emit_signal("connection_added", p_connection);
    return true;
  }
  return false;
}

bool FlexNet::remove_connection(FlexNet *p_connection) {
  if (p_connection && connections.has(p_connection)) {
    connections.erase(p_connection);
    emit_signal("connection_removed", p_connection);
    return true;
  }
  return false;
}

TypedArray<FlexNet> FlexNet::get_connections() {
  if (connections.is_empty()) {
    setup_connections();
  }
  
  TypedArray<FlexNet> result;
  for (FlexNet *conn : connections) {
    result.push_back(conn);
    print_line("pushing connection " + String(conn->get_path()));
  }

  return result;
}

void FlexNet::set_connection_paths(const TypedArray<NodePath> &p_connections) {
  connection_paths = p_connections.duplicate();
  setup_connections();
}

TypedArray<NodePath> FlexNet::get_connection_paths() const {
  return connection_paths;
}

void FlexNet::setup_connections() {
  if (!is_processing()) {
    return;
  }

  print_line("Setting up connections for " + get_path());

  Vector<FlexNet *> old_connections = connections.duplicate();
  connections.clear();

  for (FlexNet *conn : old_connections) {
    if (conn) {
      emit_signal("connection_removed", conn);
    }
  }

  for (int i = 0; i < connection_paths.size(); ++i) {
    

    FlexNet *conn = connection_paths[i].get_type() == Variant::NODE_PATH ? \
      Object::cast_to<FlexNet>(get_node_or_null(connection_paths[i])):
      nullptr;
    
    if (conn && !connections.has(conn)) {
      connections.push_back(conn);
      emit_signal("connection_added", conn);
    }
  }
}

void FlexNet::set_state(PackedInt32Array p_state) {
  for (size_t i = 0; i < _get_size_internal() && i < p_state.size(); ++i) {
    int state = p_state[i];
    if (state >= FlexNetState::V0 && state < FlexNetState::MAX) {
      states[i] = static_cast<FlexNetState>(state);
      continue; 
    }
    // Optionally, handle invalid types (e.g., skip or set to default)
    states[i] = FlexNetState::U; // Default to U if invalid
  }
}

PackedInt32Array FlexNet::get_state() const {
  PackedInt32Array result;
  for (size_t i = 0; i < _get_size_internal(); ++i) {
    result.push_back(states[i]);
  }
  return result;
}

FlexNet::FlexNet() {
  set_u(-1);
}

void FlexNet::_notification(int p_what) {
  switch (p_what) {
    case NOTIFICATION_READY:
      connections.clear();
      break;
    case NOTIFICATION_ENTER_TREE:
      // Connections may be invalid if the node was moved in the tree
      connections.clear();
      break;
    case NOTIFICATION_EXIT_TREE:
      connections.clear();
      break;
    default:
      break;
  }
}

void FlexNet::_bind_methods() {
  ClassDB::bind_method(D_METHOD("pass_state", "to"), &FlexNet::drive);
  ClassDB::bind_method(D_METHOD("set_value", "value"), &FlexNet::set_value);
  ClassDB::bind_method(D_METHOD("get_value"), &FlexNet::get_value);
  ClassDB::bind_method(D_METHOD("set_u", "mask"), &FlexNet::set_u);
  ClassDB::bind_method(D_METHOD("get_u"), &FlexNet::get_u);
  ClassDB::bind_method(D_METHOD("set_x", "mask"), &FlexNet::set_x);
  ClassDB::bind_method(D_METHOD("get_x"), &FlexNet::get_x);
  ClassDB::bind_method(D_METHOD("set_z", "mask"), &FlexNet::set_z);
  ClassDB::bind_method(D_METHOD("get_z"), &FlexNet::get_z);
  ClassDB::bind_method(D_METHOD("get_size"), &FlexNet::get_size);
  ClassDB::bind_method(D_METHOD("add_connection", "connection"), &FlexNet::add_connection);
  ClassDB::bind_method(D_METHOD("remove_connection", "connection"), &FlexNet::remove_connection);
  ClassDB::bind_method(D_METHOD("get_connections"), &FlexNet::get_connections);
  ClassDB::bind_method(D_METHOD("set_connection_paths", "connections"), &FlexNet::set_connection_paths);
  ClassDB::bind_method(D_METHOD("get_connection_paths"), &FlexNet::get_connection_paths);
  ClassDB::bind_method(D_METHOD("set_state", "state"), &FlexNet::set_state);
  ClassDB::bind_method(D_METHOD("get_state"), &FlexNet::get_state);

  ADD_PROPERTY(PropertyInfo(Variant::PACKED_INT32_ARRAY, "state", PROPERTY_HINT_ARRAY_TYPE, "int:enum/FlexNetState"), "set_state", "get_state");
  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "connection_paths", PROPERTY_HINT_ARRAY_TYPE, "Node"), "set_connection_paths", "get_connection_paths");

  ADD_SIGNAL(MethodInfo("connection_added", PropertyInfo(Variant::NODE_PATH, "connection")));
  ADD_SIGNAL(MethodInfo("connection_removed", PropertyInfo(Variant::NODE_PATH, "connection")));
}
