package ltypes

import (
	"golang.org/x/exp/constraints"
)

/*
 interfaces for circuit types
*/

type LState interface {
	any
}

type LTime interface {
	constraints.Ordered
}
