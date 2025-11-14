package llabeling

/*
 Label a struct as having some configuration of values where it will be
 treated as "null" or "empty"
*/
type optional interface {
	comparable
	IsEmpty() bool
}

/*
 Non-unique label mapping that packs elements contiguously
*/
type llabeling[T optional] []T

func (set *llabeling[T]) Add(element T) {
	for i, v := range *set {
		if v.IsEmpty() {
			(*set)[i] = element
			return
		}
	}

	*set = append(*set, element)
}

// Add an element to the desired position, growing the array if necessary
func (set *llabeling[T]) Set(element T, at int) {
	if at >= len(*set) {
		// Grow the slice to accommodate the desired position
		newSize := at + 1
		newSlice := make([]T, newSize)
		copy(newSlice, *set)
		*set = newSlice
	}
	(*set)[at] = element
}

// Get returns a pointer to the element at the specified position, or nil if out of bounds or empty
func (set *llabeling[T]) Get(at int) *T {
	if at < 0 || at >= len(*set) {
		return nil // Out of bounds
	}
	if (*set)[at].IsEmpty() {
		return nil // Empty slot
	}
	return &(*set)[at]
}

// Remove sets the element at the specified position to its "empty" value
func (set *llabeling[T]) Remove(at int) {
	if at < 0 || at >= len(*set) {
		return // Out of bounds, do nothing
	}
	var empty T // Zero value of T
	(*set)[at] = empty
}

/*
 Simple labeling scheme: A, B, C, ...
 Limit to 1-character long labels for now
*/
func Label(index int) string {
	return string(rune('A' + index))
}

func Index(label string) int {
	return int([]rune(label)[0] - 'A')
}
