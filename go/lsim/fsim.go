package lsim

import (
	c "flex-logic/collections"
	"flex-logic/lcircuit"
	"flex-logic/ltypes"
)

type Lsim[S ltypes.LState, T ltypes.LTime] struct {
	Circuit lcircuit.LCircuit[S, T]
	Events  c.LEvents[S, T]
}
