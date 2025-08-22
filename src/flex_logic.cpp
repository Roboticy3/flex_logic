#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>
#include <godot_cpp/templates/vector.hpp>

#include <godot_cpp/core/class_db.hpp>

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

  for (int i = 0; i < sizeof(int); ++i) {
    r_to->states[i] = Buckwell42[r_to->states[i] * FlexNetState::MAX + states[i]];
  }
}

void FlexNet::set_value(int p_value) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (p_value >> i)
      states[i] = FlexNetState::V1; 
    else
      states[i] = FlexNetState::V0; 
  }
}

int FlexNet::get_value() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    value |= ((states[i] & 0x01) << i);
  }
  return value;
}

void FlexNet::set_u(int p_mask) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (p_mask >> i)
      states[i] = FlexNetState::U;
  }
}

int FlexNet::get_u() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (states[i] == FlexNetState::U)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_x(int p_mask) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (p_mask >> i)
      states[i] = FlexNetState::X;
  }
}

int FlexNet::get_x() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (states[i] == FlexNetState::X)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_z(int p_mask) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if ((p_mask >> i) & 0x01)
      states[i] = FlexNetState::Z;
  }
}

int FlexNet::get_z() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (states[i] == FlexNetState::Z)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_connections(const TypedArray<NodePath> &p_connections) {
  connections.clear();
  for (int i = 0; i < p_connections.size(); ++i) {
    FlexNet *conn = p_connections[i].get_type() == Variant::NODE_PATH ? \
      Object::cast_to<FlexNet>(get_node_or_null(p_connections[i])):
      nullptr;
    connections.push_back(conn);
  }
}

TypedArray<NodePath> FlexNet::get_connections() const {
  TypedArray<NodePath> result;
  for (FlexNet *connection : connections) {
    if (connection) {
      result.push_back(get_path_to(connection));
    } else {
      result.push_back(Variant::NIL);
    }
      
  }
  return result;
}

FlexNet::FlexNet() {
  set_u(-1);
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
  ClassDB::bind_method(D_METHOD("set_connections", "connections"), &FlexNet::set_connections);
  ClassDB::bind_method(D_METHOD("get_connections"), &FlexNet::get_connections);

  ADD_PROPERTY(PropertyInfo(Variant::INT, "value"), "set_value", "get_value");
  ADD_PROPERTY(PropertyInfo(Variant::INT, "u_mask"), "set_u", "get_u");
  ADD_PROPERTY(PropertyInfo(Variant::INT, "x_mask"), "set_x", "get_x");
  ADD_PROPERTY(PropertyInfo(Variant::INT, "z_mask"), "set_z", "get_z");
  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "connections"), "set_connections", "get_connections");
}
