#include <godot_cpp/classes/node3d.hpp>
#include <godot_cpp/templates/vector.hpp>
#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>

#include <godot_cpp/core/binder_common.hpp>
#include <godot_cpp/core/class_db.hpp>

using namespace godot;

enum FlexNetState {
  V0 = 0,
  V1 = 1,
  X,
  Z,
  U,
  MAX
};

class FlexNet : public Node3D {
  GDCLASS(FlexNet, Node3D)

  FlexNetState states[sizeof(int)] = {};
  Vector<FlexNet *> connections = {};

protected:
  virtual void solver(Vector<FlexNet *> &r_event_queue);

  inline size_t _get_size_internal() const;

  static void _bind_methods();

public:

  void drive(FlexNet *r_to);

  //set each bit of p_value to the first bits of each element in state
  void set_value(int p_value);
  int get_value() const;

  //convert all high bits to special values
  void set_u(int p_mask);
  int get_u() const;

  void set_x(int p_mask);
  int get_x() const;

  void set_z(int p_mask);
  int get_z() const;

  inline size_t get_size() const;

  void set_connections(const TypedArray<NodePath> &p_connections);
  TypedArray<NodePath> get_connections() const;

  void set_state(PackedInt32Array p_state);
  PackedInt32Array get_state() const;

  FlexNet();
};