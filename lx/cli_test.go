package lx

import (
	"testing"
)

// Basic sanity test to ensure NewCommand constructs a command without panicking
// and that it has the expected name.
func TestNewCommand_Basic(t *testing.T) {
	cmd := NewCommand()
	if cmd == nil {
		t.Fatal("NewCommand() returned nil")
	}
	if cmd.Name != "lx" {
		t.Fatalf("NewCommand().Name = %q, want %q", cmd.Name, "lx")
	}
}

