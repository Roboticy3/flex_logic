package lcircuit

import (
	"golang.org/x/exp/constraints"
)

/*
 interfaces for circuit types
*/

type lstate interface {
	any
}

type ltime interface {
	constraints.Ordered
}
