package Lcircuit

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

func (set *Llabeling[T]) Add(element T, start int64) int64 {
	if start >= 0 && start < int64(len(*set)) {
		for i, v := range (*set)[start:] {
			if v.IsEmpty() {
				(*set)[i] = element
				return int64(i)
			}
		}
	}

	*set = append(*set, element)
	return int64(len(*set) - 1)
}

// Add an element to the desired position, growing the array if necessary
func (set *Llabeling[T]) Set(element T, at int64) {
	if at >= int64(len(*set)) {
		// Grow the slice to accommodate the desired position
		newSize := at + 1
		newSlice := make([]T, newSize)
		copy(newSlice, *set)
		*set = newSlice
	}
	(*set)[at] = element
}

// Get returns a pointer to the element at the specified position, or nil if out of bounds or empty
func (set *Llabeling[T]) Get(at int64) *T {
	if at < 0 || at >= int64(len(*set)) {
		return nil // Out of bounds
	}
	if (*set)[at].IsEmpty() {
		return nil // Empty slot
	}
	return &(*set)[at]
}

// Remove sets the element at the specified position to its "empty" value
func (set *Llabeling[T]) Remove(at int64) {
	if at < 0 || at >= int64(len(*set)) {
		return // Out of bounds, do nothing
	}
	var empty T // Zero value of T
	(*set)[at] = empty
}

/*
Simple labeling scheme: A, B, C, ...
Limit to 1-character long labels for now
*/
func Label(index int64) string {
	return strconv.FormatInt(index, 26)
}

func Index(label string) (int64, error) {
	return strconv.ParseInt(label, 26, 64)
}

// Useful example for testing
type string_component string

func (s string_component) IsEmpty() bool {
	return s == ""
}

func (s string_component) GetEmpty() string_component {
	return ""
}
