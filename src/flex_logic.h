#include <godot_cpp/classes/node.hpp>
#include <godot_cpp/templates/vector.hpp>
#include <godot_cpp/variant/array.hpp>

using namespace godot;


class FlexNet : public Node {
  GDCLASS(FlexNet, Node)

  //array of bytes
  //the first bit of each byte concatenates into an int
  //the second and third bits encode special values
  //u (unset): 0x01*
  //x (conflict): 0x10*
  //z: 0x11*
  int8_t states[sizeof(int)] = { 0 };
  Vector<FlexNet *> connections = {};

protected:
  //hood virtual
  static void _default_solver(const FlexNet &net, List<FlexNet *> &r_event_queue);
  void (*solver)(const FlexNet &, List<FlexNet *> &) = _default_solver;

  static void _bind_methods();

public:

  void pass_state(FlexNet *r_to);

  //each bit of `p_value` corresponds to the value bit of an element in `states` 
  //the other setters work the same way.
  void set_value(int p_value);
  int get_value() const;

  void set_u(int p_mask);
  int get_u() const;

  void set_x(int p_mask);
  int get_x() const;

  void set_z(int p_mask);
  int get_z() const;

  void set_connections(const Array &p_connections);
  Array get_connections() const;

  FlexNet();
};