#pragma once

#include <cstdint>

#define WIDTH sizeof(int)<<3

enum WireState {
  V0 = 0,
  V1 = 1,
  X,
  Z,
  U,
  MAX
};

struct FlexNetState {
  WireState states[WIDTH];
  uint16_t solver;
  std::vector<FlexNetState*> connections = {};
} typedef FlexNetState;