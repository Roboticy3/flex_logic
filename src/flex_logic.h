#pragma once

#include <vector>
#include <queue>
#include <unordered_map>
#include <functional>

#include <godot_cpp/classes/node3d.hpp>
#include <godot_cpp/templates/vector.hpp>
#include <godot_cpp/variant/array.hpp>
#include <godot_cpp/variant/node_path.hpp>

#include <godot_cpp/core/binder_common.hpp>
#include <godot_cpp/core/class_db.hpp>

#include "flex_state.h"
#include "flex_net.h"

using namespace godot;

struct FlexConnection {
  const FlexNet* net = nullptr;
};

class FlexLogic : public Node3D {
  GDCLASS(FlexLogic, Node3D)
  
  /**
   * Keep good locality on the simulation state by separating it from FlexNet.
   */
  std::vector<FlexNetState> net_states = {};

  /**
   * Keep track of nets in the same order as `net_states`. Use Vector for easier
   * search.
   */
  Vector<const FlexNet *> nets = {};

  size_t get_net_count() const;

  void restore_connections_from(size_t p_id);

  /**
   * Each gate type has a "solver" function which takes a state and returns
   * events. Events are added to a total queue and processed in order.
   * 
   * The order of `const vector<FlexNetState *> &` is the same as the associated
   * element in `connections`, which should be passed as a direct reference.
   */
  std::vector<std::function<void(const std::vector<FlexNetState *> &, std::queue<size_t>)>> solvers = {
    FlexLogic::solver_wire
  };

  std::queue<size_t> event_queue = {};

  protected:
    static void _bind_methods();

  public:
    static void drive(const FlexNetState &p_from, FlexNetState &p_to);
    static void solver_wire(const std::vector<FlexNetState *> &p_states, std::queue<size_t> &r_event_queue);

    void restore_connections();

    /**
     * Push all events caused by `p_start_index`, then stop.
     */
    bool step_from(const FlexNet *p_start_net);

    bool add_net(const FlexNet *p_net, PackedInt32Array p_initial_state = PackedInt32Array());
    bool remove_net(const FlexNet *p_net);

    Array get_nets() const;
    TypedArray<PackedInt32Array> get_state() const;

    FlexLogic();
  
};