#pragma once

#include <functional>
#include <queue>

#include <godot_cpp/classes/node3d.hpp>
#include <godot_cpp/templates/vector.hpp>
#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>

#include <godot_cpp/core/binder_common.hpp>
#include <godot_cpp/core/class_db.hpp>

#include "flex_state.h"
#include "flex_net.h"

using namespace godot;

class FlexLogic : public Node3D {
  GDCLASS(FlexLogic, Node3D)
  
  /**
   * Keep good locality on the simulation state by separating it from FlexNet.
   */
  FlexNetState* net_states = {};

  /**
   * Keep track of nets in the same order as `net_states`
   */
  Vector<Vector<FlexNetState *>> connections;

  /**
   * Keep track of nets in the same order as `net_states`
   */
  Vector<const FlexNet *> nets = {};

  size_t get_net_count() const;

  /**
   * Each gate type has a "solver" function which takes a state and returns
   * events. Events are added to a total queue and processed in order.
   * 
   * The order of `const Vector<FlexNetState *> &` is the same as the associated
   * element in `connections`, which should be passed as a direct reference.
   */
  Vector<std::function<void(const Vector<FlexNetState *> &, std::queue<size_t>)>> solvers;

  std::queue<size_t> event_queue = {};

  protected:
    static void _bind_methods();

  public:
    static void drive(const FlexNetState &p_from, FlexNetState &p_to);
    static void solver_wire(const Vector<FlexNetState *> &p_states, std::queue<size_t> &r_event_queue);

    /**
     * Push all events caused by `p_start_index`, then stop.
     */
    void step_from(size_t p_start_index);

    bool add_net(const FlexNet *p_net);
    bool remove_net(const FlexNet *p_net);

    FlexLogic();
  
};