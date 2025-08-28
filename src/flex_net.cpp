#include <cstdint>

#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>

#include <godot_cpp/core/binder_common.hpp>
#include <godot_cpp/core/class_db.hpp>

#include "flex_state.h"
#include "flex_net.h"
#include "flex_logic.h"

using namespace godot;

size_t FlexNet::get_size() const {
  return WIDTH;
}

NodePath FlexNet::get_logic_path() const {
  return logic_path;
}

void FlexNet::set_logic_path(const NodePath &p_path) {
  if (logic_path == p_path) {
    return;
  }

  FlexLogic *new_logic = p_path.is_empty() ? nullptr : \
    Object::cast_to<FlexLogic>(get_node_or_null(p_path));
  
  set_logic(new_logic);
  update_logic_path();
}

void FlexNet::set_logic(FlexLogic *p_logic) {
  if (logic == p_logic) {
    return;
  }

  if (logic) {
    logic->remove_net(this);
  }

  logic = p_logic;

  if (logic) {
    logic->add_net(this);
    logic->restore_connections();
    print_line("FlexNet::set_logic: added net to logic " + String(logic->get_path()));
  } else {
    //print_line("FlexNet::set_logic: logic is null, not adding net");
  }
}

void FlexNet::update_logic_path() {
  if (logic) {
    logic_path = get_path_to(logic);
  } else {
    logic_path = NodePath();
  }
}

FlexLogic *FlexNet::get_logic() const {
  return logic;
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
    //print_line("pushing connection " + String(conn->get_path()));
  }

  return result;
}

Vector<FlexNet *> FlexNet::get_connections_raw() const {
  return connections;
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

  //print_line("Setting up connections for " + get_path());

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

PackedInt32Array FlexNet::get_state() const {
  return PackedInt32Array(state);
}

void FlexNet::set_state(PackedInt32Array p_state) {
  if (p_state.size() != WIDTH) {
    ERR_PRINT("FlexNet::set_state: Invalid state size " + itos(p_state.size()) + ", expected " + itos(WIDTH));
    return;
  }

  for (int i = 0; i < WIDTH; ++i) {
    int val = p_state[i];
    if (val < 0 || val >= WireState::MAX) {
      ERR_PRINT("FlexNet::set_state: Invalid state value " + itos(val) + " at index " + itos(i));
      state[i] = (int)WireState::U;
      continue;
    }
    state[i] = val;
  }

  emit_signal("state_changed", get_state());
}

void FlexNet::set_solver(int p_solver) {
  if (p_solver < 0 || p_solver >= UINT16_MAX) {
    ERR_PRINT("FlexNet::set_solver: Invalid solver " + itos(p_solver));
    return;
  }
  solver = (uint16_t)p_solver;
}

uint16_t FlexNet::get_solver() const {
  return solver;
}

void FlexNet::_notification(int p_what) {
  switch (p_what) {
    case NOTIFICATION_READY:
      set_logic_path(logic_path);
      connections.clear();
      break;
    case NOTIFICATION_ENTER_TREE:
      // Connections may be invalid if the node was moved in the tree
      set_logic_path(logic_path);
      connections.clear();
      break;
    case NOTIFICATION_EXIT_TREE:
      connections.clear();
      logic = nullptr;
      update_logic_path();
      break;
    default:
      break;
  }
}

FlexNet::FlexNet() {
  state.resize(WIDTH);
  for (int i = 0; i < WIDTH; ++i) {
    state[i] = (int)WireState::U;
  }
}

void FlexNet::_bind_methods() {
  ClassDB::bind_method(D_METHOD("get_size"), &FlexNet::get_size);
  ClassDB::bind_method(D_METHOD("add_connection", "connection"), &FlexNet::add_connection);
  ClassDB::bind_method(D_METHOD("remove_connection", "connection"), &FlexNet::remove_connection);
  ClassDB::bind_method(D_METHOD("get_connections"), &FlexNet::get_connections);
  ClassDB::bind_method(D_METHOD("get_logic"), &FlexNet::get_logic);
  ClassDB::bind_method(D_METHOD("set_logic_path", "path"), &FlexNet::set_logic_path);
  ClassDB::bind_method(D_METHOD("get_logic_path"), &FlexNet::get_logic_path);
  ClassDB::bind_method(D_METHOD("set_connection_paths", "connections"), &FlexNet::set_connection_paths);
  ClassDB::bind_method(D_METHOD("get_connection_paths"), &FlexNet::get_connection_paths);
  ClassDB::bind_method(D_METHOD("set_state", "state"), &FlexNet::set_state);
  ClassDB::bind_method(D_METHOD("get_state"), &FlexNet::get_state);
  ClassDB::bind_method(D_METHOD("set_solver", "solver"), &FlexNet::set_solver);
  ClassDB::bind_method(D_METHOD("get_solver"), &FlexNet::get_solver);

  ADD_PROPERTY(PropertyInfo(Variant::NODE_PATH, "logic_path", PROPERTY_HINT_NODE_TYPE, "Node"), "set_logic_path", "get_logic_path");
  ADD_PROPERTY(PropertyInfo(Variant::ARRAY, "connection_paths", PROPERTY_HINT_ARRAY_TYPE, "Node"), "set_connection_paths", "get_connection_paths");
  ADD_PROPERTY(PropertyInfo(Variant::PACKED_INT32_ARRAY, "state", PROPERTY_HINT_ARRAY_TYPE, "int"), "set_state", "get_state");
  ADD_PROPERTY(PropertyInfo(Variant::INT, "solver"), "set_solver", "get_solver");

  ADD_SIGNAL(MethodInfo("connection_added", PropertyInfo(Variant::NODE_PATH, "connection")));
  ADD_SIGNAL(MethodInfo("connection_removed", PropertyInfo(Variant::NODE_PATH, "connection")));
  ADD_SIGNAL(MethodInfo("state_changed", PropertyInfo(Variant::PACKED_INT32_ARRAY, "state")));

}