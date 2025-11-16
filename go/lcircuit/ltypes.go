package lcircuit

import (
	"golang.org/x/exp/constraints"
)

/*
 interfaces for circuit types
*/

type Lstate interface {
	any
}

type Ltime interface {
	constraints.Ordered
}
