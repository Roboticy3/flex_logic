#include <godot_cpp/classes/node.hpp>
#include <godot_cpp/templates/vector.hpp>
#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>

#include <godot_cpp/core/binder_common.hpp>

using namespace godot;

enum FlexNetState {
  V0 = 0,
  V1 = 1,
  X,
  Z,
  U,
  MAX
};

class FlexNet : public Node {
  GDCLASS(FlexNet, Node)

  FlexNetState states[sizeof(int)] = {};
  Vector<FlexNet *> connections = {};

protected:
  virtual void solver(Vector<FlexNet *> &r_event_queue);

  inline void set_special(int p_mask, FlexNetState p_value);
  inline int get_special(int p_mask, FlexNetState p_value) const;

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

  void set_connections(const TypedArray<NodePath> &p_connections);
  TypedArray<NodePath> get_connections() const;

  FlexNet();
};