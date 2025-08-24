#pragma once

#include <godot_cpp/classes/node3d.hpp>
#include <godot_cpp/variant/node_path.hpp>
#include <godot_cpp/variant/packed_int32_array.hpp>

#include <godot_cpp/core/binder_common.hpp>
#include <godot_cpp/core/class_db.hpp>

#include "flex_state.h"

class FlexLogic;

using namespace godot;

class FlexNet : public Node3D {
  GDCLASS(FlexNet, Node3D)

  FlexLogic *logic = nullptr;
  NodePath logic_path = NodePath();

  void update_logic_path();

  Vector<FlexNet *> connections = {};
  TypedArray<NodePath> connection_paths = TypedArray<NodePath>();

  /**
   * Load `connections` from `connection_paths`.
   * Each invalid path is skipped.
   */
  void setup_connections();

  /**
   * Store the starting state of the net. Read-only from scripts.
   */
  PackedInt32Array state = PackedInt32Array();

  /**
   * The index of this net in `logic->solvers`
   */
  uint16_t solver = 0;

protected:
  static void _bind_methods();
  void _notification(int p_what);

public:
  size_t get_size() const;

  void set_logic_path(const NodePath &p_path);
  NodePath get_logic_path() const;

  void set_logic(FlexLogic *p_logic);
  FlexLogic *get_logic() const;

  bool add_connection(FlexNet *p_connection);
  bool remove_connection(FlexNet *p_connection);
  TypedArray<FlexNet> get_connections();
  Vector<FlexNet *> get_connections_raw() const;

  void set_connection_paths(const TypedArray<NodePath> &p_connections);
  TypedArray<NodePath> get_connection_paths() const;

  void set_state(PackedInt32Array p_state);
  PackedInt32Array get_state() const;

  void set_solver(int p_solver);
  uint16_t get_solver() const;

  FlexNet();
};