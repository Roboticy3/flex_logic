package lcircuit

import (
	"strconv"
)

/*
Label a struct as having some configuration of values where it will be
treated as "null" or "empty"
*/
type optional interface {
	IsEmpty() bool
}

/*
Non-unique label mapping that packs elements contiguously
*/
type Llabeling[T optional] []T

func (set *Llabeling[T]) Len() int {
	return len(*set)
}

func (set *Llabeling[T]) Add(element T, start int) int {
	if start >= 0 && start < int(len(*set)) {
		for i, v := range (*set)[start:] {
			if v.IsEmpty() {
				(*set)[start+i] = element
				return int(start + i)
			}
		}
	}

	*set = append(*set, element)
	return int(len(*set) - 1)
}

// Add an element to the desired position, growing the array if necessary
func (set *Llabeling[T]) Set(element T, at int) {
	if at >= int(len(*set)) {
		// Grow the slice to accommodate the desired position
		newSize := at + 1
		newSlice := make([]T, newSize)
		copy(newSlice, *set)
		*set = newSlice
	}
	(*set)[at] = element
}

// Get returns a pointer to the element at the specified position, or nil if out of bounds or empty
func (set *Llabeling[T]) Get(at int) *T {
	if at < 0 || at >= int(len(*set)) {
		return nil // Out of bounds
	}
	if (*set)[at].IsEmpty() {
		return nil // Empty slot
	}
	return &(*set)[at]
}

// Remove sets the element at the specified position to a value that returns true on empty.IsEmpty()
func (set *Llabeling[T]) Remove(at int, empty T) {
	if at < 0 || at >= int(len(*set)) {
		return // Out of bounds, do nothing
	}
	(*set)[at] = empty
}

/*
Simple labeling scheme: A, B, C, ...
Limit to 1-character long labels for now
*/
func Label(index int) string {
	return strconv.FormatInt(int64(index), 26)
}

func Index(label string) (int, error) {
	i, err := strconv.ParseInt(label, 26, 64)
	return int(i), err
}

// Useful example for testing
type string_component string

func (s string_component) IsEmpty() bool {
	return s == ""
}

func (s string_component) GetEmpty() string_component {
	return ""
}
