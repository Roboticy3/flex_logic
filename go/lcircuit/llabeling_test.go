package lcircuit

import "testing"

func TestAdd(t *testing.T) {
	var labels Llabeling[string_label]
	labels.Add("example", 0)
	labels.Add("test", 0)

	if len(labels) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(labels))
	}
}

func TestSet(t *testing.T) {
	var labels Llabeling[string_label]
	labels.Set("example", 2)

	if len(labels) != 3 {
		t.Errorf("Expected length 3, got %d", len(labels))
	}
	if labels[2] != "example" {
		t.Errorf("Expected 'example' at position 2, got %v", labels[2])
	}
}

func TestGet(t *testing.T) {
	var labels Llabeling[string_label]
	labels.Set("example", 1)

	if labels.Get(1) == nil || *labels.Get(1) != "example" {
		t.Errorf("Expected 'example' at position 1, got nil or incorrect value")
	}
	if labels.Get(0) != nil {
		t.Errorf("Expected nil at position 0, got %v", labels.Get(0))
	}
}

func TestRemove(t *testing.T) {
	var labels Llabeling[string_label]
	labels.Set("example", 1)
	labels.Remove(1, "")

	if labels.Get(1) != nil {
		t.Errorf("Expected nil at position 1 after removal, got %v", labels.Get(1))
	}
}
