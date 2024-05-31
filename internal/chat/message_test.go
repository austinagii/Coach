package chat

import (
	"testing"
)

func TestString(t *testing.T) {
	m := NewUserMessage("test")
	expected := "user: test"
	actual := m.String()
	if actual != expected {
		t.Fatalf("Invalid value; expected: '%s', got '%s'", expected, actual)
	}

}
