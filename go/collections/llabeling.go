package collections

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
type LLabeling[T optional] []T

func (set *LLabeling[T]) Len() int {
	return len(*set)
}

func (set *LLabeling[T]) Add(element T, start int) Label {
	if start >= 0 && start < len(*set) {
		for i, v := range (*set)[start:] {
			if v.IsEmpty() {
				(*set)[start+i] = element
				return Label(start + i)
			}
		}
	}

	*set = append(*set, element)
	return Label(len(*set) - 1)
}

// Add an element to the desired position, growing the array if necessary
func (set *LLabeling[T]) Set(element T, at Label) {
	if at >= Label(len(*set)) {
		// Grow the slice to accommodate the desired position
		newSize := at + 1
		newSlice := make([]T, newSize)
		copy(newSlice, *set)
		*set = newSlice
	}
	(*set)[at] = element
}

// Get returns a poLabeler to the element at the specified position, or nil if out of bounds or empty
func (set *LLabeling[T]) Get(at Label) *T {
	if at < 0 || at >= Label(len(*set)) {
		return nil // Out of bounds
	}
	if (*set)[at].IsEmpty() {
		return nil // Empty slot
	}
	return &(*set)[at]
}

// Remove sets the element at the specified position to a value that returns true on empty.IsEmpty()
func (set *LLabeling[T]) Remove(at Label, empty T) {
	if at < 0 || at >= Label(len(*set)) {
		return // Out of bounds, do nothing
	}
	(*set)[at] = empty
}

func (set *LLabeling[T]) Compress() []T {
	result := []T{}
	for _, item := range *set {
		if !item.IsEmpty() {
			result = append(result, item)
		}
	}

	return result
}

/*
Simple labeling scheme: A, B, C, ...
Limit to 1-character long labels for now
*/
func Name(index Label) string {
	return strconv.FormatInt(int64(index), 26)
}

func Index(name string) (Label, error) {
	i, err := strconv.ParseInt(name, 26, 64)
	return Label(i), err
}

// Example implementaations
type Label int

// Empty value can be referenced to create empty entries
const LABEL_EMPTY = Label(-1)

// IsEmpty reduces to a comparison with the empty value
func (i Label) IsEmpty() bool {
	return i == LABEL_EMPTY
}

// Second example
type StringLabel string

const STRING_EMPTY = ""

func (s StringLabel) IsEmpty() bool {
	return s == STRING_EMPTY
}
