package util

import (
	"testing"
)

func TestNewId(t *testing.T) {
	// Reset current for test isolation
	current = 1
	prefix := "test"
	id1 := NewUniqueId(prefix)
	id2 := NewUniqueId(prefix)
	if id1 == id2 {
		t.Errorf("Expected unique ids, got %v and %v", id1, id2)
	}
	if id1 == "" || id2 == "" {
		t.Error("Expected non-empty ids")
	}
}

func TestNewUniqeTypeId(t *testing.T) {
	type TestType struct{}
	id1 := NewUniqeTypeId[TestType]()
	id2 := NewUniqeTypeId[TestType]()
	if id1 == id2 {
		t.Errorf("Expected unique type ids, got %v and %v", id1, id2)
	}
	if id1 == "" || id2 == "" {
		t.Error("Expected non-empty unique type ids")
	}
}

func TestBroadcast(t *testing.T) {
	msg := "hello"
	cmd := Broadcast(msg)
	result := cmd()
	if result != msg {
		t.Errorf("Expected %v, got %v", msg, result)
	}
}
