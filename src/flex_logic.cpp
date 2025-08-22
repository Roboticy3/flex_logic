#include "flex_logic.h"

using namespace godot;

void FlexNet::_default_solver(const FlexNet &p_net, List<FlexNet *> &r_event_queue) {
  //Copy our value into all connections, then add them to the event queue
  for (FlexNet *connection : p_net.connections) {
    connection->set_value(p_net.get_value());
    r_event_queue.push_back(connection);
  } 
}

void FlexNet::pass_state(FlexNet *r_to) {
  memcpy(r_to->states, states, sizeof(states));
}

void FlexNet::set_value(int p_value) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    // Set the value bit (bit 0) for each state byte
    states[i] = (states[i] & ~0x01) | ((p_value >> i) & 0x01);
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
    if ((p_mask >> i) & 0x01)
      states[i] = (states[i] & ~0x02) | 0x02; // Set bit 1
    else
      states[i] &= ~0x02; // Clear bit 1
  }
}

int FlexNet::get_u() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (states[i] & 0x02)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_x(int p_mask) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if ((p_mask >> i) & 0x01)
      states[i] = (states[i] & ~0x04) | 0x04; // Set bit 2
    else
      states[i] &= ~0x04; // Clear bit 2
  }
}

int FlexNet::get_x() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if (states[i] & 0x04)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_z(int p_mask) {
  for (size_t i = 0; i < sizeof(int); ++i) {
    if ((p_mask >> i) & 0x01)
      states[i] = (states[i] & ~0x06) | 0x06; // Set bits 1 and 2
    else
      states[i] &= ~0x06; // Clear bits 1 and 2
  }
}

int FlexNet::get_z() const {
  int value = 0;
  for (size_t i = 0; i < sizeof(int); ++i) {
    if ((states[i] & 0x06) == 0x06)
      value |= (1 << i);
  }
  return value;
}

void FlexNet::set_connections(const Array &p_connections) {
  connections.clear();
  for (int i = 0; i < p_connections.size(); ++i) {
    FlexNet *conn = Object::cast_to<FlexNet>(p_connections[i].get_validated_object());
    connections.push_back(conn);
  }
}

Array FlexNet::get_connections() const {
  Array result;
  for (FlexNet *connection : connections) {
    Variant v = Variant(connection);
    result.push_back(v);
  }
  return result;
}

FlexNet::FlexNet() {}

void FlexNet::_bind_methods() {
  ClassDB::bind_method(D_METHOD("pass_state", "to"), &FlexNet::pass_state);
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
  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "connections", PROPERTY_HINT_ARRAY_TYPE, "Node"), "set_connections", "get_connections");
}
